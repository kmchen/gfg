package elasticsearch

import (
	"context"
	"testing"
)

const esUrl = "http://localhost:9200"

func TestSearchQuery(t *testing.T) {
	ctx := context.Background()
	esClient, err := NewElasticSearchClient(esUrl, ctx)
	if err != nil {
		t.Errorf("Init elastic search client gets an error")
	}
	var query = SearchQuery{
		Search: Param{
			Key:   "play_name",
			Value: "Henry IV",
		},
		Pagination: PaginationParam{
			Limit:  3,
			Offset: 0,
		},
		Sort: SortParam{
			Key: "line_id",
			Asc: false,
		},
		Filter: Param{
			Key:   "speaker",
			Value: "KING HENRY IV",
		},
	}
	searchService := esClient.SearchService(query)
	customers, err := esClient.DoSearch(searchService, ctx)
	if err != nil {
		t.Errorf("Search customer by id gets an error %v", err)
	}
	if len(customers) <= 0 {
		t.Errorf("Found 0 customers")
	}
}
