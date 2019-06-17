package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	auth "github.com/gfg/authentication"
	es "github.com/gfg/elasticsearch"

	"github.com/gorilla/mux"
)

const esUrl = "http://es01:9200"

const index = "shakespeare"

func populate() {
	cmd := exec.Command("sh", "-c", "./populate.sh")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	log.Println(" =========== Import start ===========")
	log.Println(out.String())
	if err != nil {
		log.Println("Fail to import data to elasticsearch")
		log.Panic(err)
	}
	log.Println(" =========== Import finished ===========")
}

func main() {
	ctx := context.Background()

	// Init elasticsearch client
	esClient, err := es.NewElasticSearchClient(esUrl, ctx)
	if err != nil {
		log.Fatal("Fail to initiate elasticsearch client")
	}

	// Check if data is imported otherwise populate data
	if !esClient.IsDataImported(index, ctx) {
		populate()
	}

	authentication := &auth.Authentication{}
	authentication.Populate()

	// Construct routes v1 and v2 return same data
	router := mux.NewRouter()
	var api = router.PathPrefix("/shakespeare").Subrouter()
	var api1 = api.PathPrefix("/v1").Subrouter()
	api1.HandleFunc("", authentication.Middleware(esClient.Handler))
	// V2 doesn't require authentication
	var api2 = api.PathPrefix("/v2").Subrouter()
	api2.HandleFunc("", esClient.Handler)

	router.HandleFunc("/token", authentication.TokenRequestHandler)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("============ Starting Server ============")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
