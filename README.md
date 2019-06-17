### Endpoints

 - POST /token
        - Aquire new token for authentication
    ```sh
    $ curl -d '{"Id":"clientId0"}' -H "Content-Type: application/json" -X POST http://localhost:8000/token
    $ eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRJZCI6ImNsaWVudElkMCIsImV4cCI6IjIwMTktMDYtMTdUMDg6MjY6MDQuODE1MjM1M1oifQ.6cBVFrnneBabyxdMh8i5VbRPiSXzeG1EKOF3KipMpLM
    ```
 - GET /shakespeare/v1
        - Get Shakespeare data with authentication
     ```sh
    $ curl -H "X-Session-Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRJZCI6ImNsaWVudElkMCIsImV4cCI6IjIwMTktMDYtMTdUMDg6MjY6MDQuODE1MjM1M1oifQ.6cBVFrnneBabyxdMh8i5VbRPiSXzeG1EKOF3KipMpLM" http://localhost:8000/shakespeare/v1?sort_by=line_id&offset=0&limit=10&filter=speaker:KING%20HENRY%20IV&play_name=Henry%20IV
    ```
- GET /shakespeare/v2
        - Get Shakespeare data with without authentication
     ```sh
    $ curl http://localhost:8000/shakespeare/v2?sort_by=line_id&offset=0&limit=10&filter=speaker:KING%20HENRY%20IV&play_name=Henry%20IV
    ```

### Summary
- Serving Shakespeare story data and can be searched by play name, filtered with speaker and sorted by line id
- The server will automatically import Shakespeare data into ElasticSearch for the first time (~5mins). The server will be ready until you see in the log
============ Starting Server ============

### Installation
 ```sh
$ git clone https://github.com/kmchen/gfg.git
$ docker-compose up
```

### TODO
- The search criteria should be more elatorated. It's now limited play_name only.
- Need more more overall unit testings
- In order to run the elasticsearch test, the elastic containers must be up and running. Need to change this.
- Should use net/http/httptest to better test http server and middleware
