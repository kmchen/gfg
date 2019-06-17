package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
)

const playName = "play_name"
const index = "shakespeare"
const indexNotFoundError = "index_not_found_exception"
const retryInterval = 10 * time.Second
const retry = 6

type ShakespeareData struct {
	Type          string      `json:"type"`
	Line_id       int         `json:"line_id"`
	Play_name     string      `json:"play_name"`
	Speech_number interface{} `json:"speech_number"`
	Line_number   string      `json:"line_number"`
	Speaker       string      `json:"speaker"`
	Text_entry    string      `json:"text_entry"`
}

type Param struct {
	Key   string
	Value string
}

type SortParam struct {
	Key string
	Asc bool
}

type PaginationParam struct {
	Limit  int
	Offset int
}

type SearchQuery struct {
	Search     Param
	Pagination PaginationParam
	Sort       SortParam
	Filter     Param
}

type ElasticSearchClient struct {
	client *elastic.Client
}

func NewElasticSearchClient(url string, ctx context.Context) (*ElasticSearchClient, error) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		log.Println("Fail to initiate elastic client: Retry in 10 seconds for 6 times")
		for i := 0; i < retry; i += 1 {
			timer := time.NewTimer(retryInterval)
			<-timer.C
			log.Printf("================= %v retry to connect to ElasticSearch ================", i+1)
			client, err = elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
			if err != nil {
				log.Println("Fail to initiate elastic client: Retry in 10 seconds for 6 times")
			} else {
				break
			}
		}
		if err != nil {
			return nil, fmt.Errorf("Fail to initiate elastic client")
		}
	}
	_, _, err = client.Ping(url).Do(ctx)
	if err != nil {
		log.Println("Fail to ping elastic search cluster")
		return nil, err
	}

	return &ElasticSearchClient{client}, nil
}

// SearchService is for testing purse only
func (esClient *ElasticSearchClient) SearchService(query SearchQuery) *elastic.SearchService {
	searchQuery := elastic.NewBoolQuery()
	matchQuery1 := elastic.NewMatchQuery(query.Search.Key, query.Search.Value)
	matchQuery2 := elastic.NewMatchQuery(query.Filter.Key, query.Filter.Value)
	searchQuery = searchQuery.Must(matchQuery1).Filter(matchQuery2)
	return esClient.client.Search().
		Index(index).
		Query(searchQuery).
		From(query.Pagination.Offset).
		Size(query.Pagination.Limit).
		Sort(query.Sort.Key, query.Sort.Asc).
		Pretty(true)
}

func (esClient *ElasticSearchClient) DoSearch(searchService *elastic.SearchService, ctx context.Context) ([]ShakespeareData, error) {
	searchResult, err := searchService.Do(ctx)
	if err != nil {
		log.Printf("Fail to search customer with query\n")
		return nil, err
	}

	if searchResult.TotalHits() > 0 {
		log.Printf("Found a total of %d data\n", searchResult.TotalHits())
		var hitsLength = len(searchResult.Hits.Hits)
		var dataList = make([]ShakespeareData, hitsLength)
		for index, hit := range searchResult.Hits.Hits {
			var data ShakespeareData
			err := json.Unmarshal(hit.Source, &data)
			if err != nil {
				log.Printf("Fail to unmarshal search result %v\n", string(hit.Source))
				return nil, err
			}
			dataList[index] = data
		}
		return dataList, nil
	}
	return nil, nil
}

func (esClient *ElasticSearchClient) IsDataImported(index string, ctx context.Context) bool {
	_, err := esClient.client.Search().Index(index).Do(ctx)
	e, _ := err.(*elastic.Error)
	if err != nil && e.Details.Type == indexNotFoundError {
		return false
	}
	return true
}

func (esClient *ElasticSearchClient) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	query := r.URL.Query()

	searchService := esClient.client.Search()
	searchQuery := elastic.NewBoolQuery()
	var searchMatchQuery *elastic.MatchQuery
	var filterMatchQuery *elastic.MatchQuery

	play_name := query.Get(playName)
	if play_name != "" {
		searchMatchQuery = elastic.NewMatchQuery(playName, play_name)
		searchQuery = searchQuery.Must(searchMatchQuery)
	}
	sort_by := query.Get("sort_by")
	if sort_by != "" {
		searchService = searchService.Sort(sort_by, true)
	}
	limit := query.Get("limit")
	if limit != "" {
		limitInt, _ := strconv.Atoi(limit)
		searchService = searchService.Size(limitInt)
	}
	offset := query.Get("offset")
	if offset != "" {
		offsetInt, _ := strconv.Atoi(offset)
		searchService = searchService.From(offsetInt)
	}
	filter := query.Get("filter")
	if filter != "" {
		filters := strings.Split(filter, ":")
		filterMatchQuery = elastic.NewMatchQuery(filters[0], filters[1])
		searchQuery = searchQuery.Filter(filterMatchQuery)
	}
	searchService = searchService.Index(index).Query(searchQuery).Pretty(true)
	searchResult, err := esClient.DoSearch(searchService, ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonCustomers, err := json.MarshalIndent(searchResult, " ", " ")
	if err != nil {
		log.Printf("Fail to marshal search data %v\n", searchResult)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonCustomers)
}
