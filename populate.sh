curl -H 'Content-Type: application/x-ndjson' -XPOST 'es01:9200/shakespeare/doc/_bulk?pretty' --data-binary @shakespeare_6.0.json
