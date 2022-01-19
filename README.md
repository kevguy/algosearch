# AlgoSearch

## Introduction


## Usage

### Prerequisities

Make sure you have a CouchDB database set up and a working Algorand node, ideally an archival node.

If you want to do tracing, you can set up Zipkin too.

### How to run

#### Using Docker

Run the following command to start the container which already contains all the three services (RESTful API, metrics and frontend):

```sh
docker run \
  -e ALGOSEARCH_WEB_ENABLE_SYNC=true \
  -e ALGOSEARCH_WEB_SYNC_INTERNAL=5s \
  -e ALGOSEARCH_COUCH_DB_HOST=234.567.89.0:5984 \
  -e ALGOSEARCH_COUCH_DB_USER=algorand \
  -e ALGOSEARCH_COUCH_DB_PASSWORD=algorand \
  -e ALGOSEARCH_ALGOD_PROTOCOL=http \
  -e ALGOSEARCH_ALGOD_ADDR=234.567.89.0:4001 \
  -e ALGOSEARCH_ALGOD_TOKEN=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa \
  -e ALGOSEARCH_ZIPKIN_REPORTER_URI=http://234.567.89.0:9411/api/v2/spans \
  -e NEXT_PUBLIC_API_URL=http://0.0.0.0:5000 \
  -e METRICS_COLLECT_FROM=http://0.0.0.0:4000/debug/vars \
  algosearch:1.1
```

#### Using Docker-Compose

Go inside `PROJECT/zard/compose/compose-config.yaml` and change the environment variables accordingly for `algosearch-backend`. Then run the following command to start `docker-compose`:

```sh
make up
```

Use this command to see the logs:

```sh
make logs
```

Use this command to stop the containers:

```sh
make down
```

### Local

To run AlgoSearch locally, you need to have the following things:

- npm/yarn, for building and starting the frontend service
- golang, for building and starting the backend services
- a couchdb connection, for the backend RESTful APi to store and retrieve data

### Backend

Both the restful API and metric services are configurable. Run the following commands to see what variables that can be configured through command line arguments or environment variables:

```sh

# RESTful API
go run ./backend/app/algosearch/main.go --help

# Metrics
go run ./backend/app/sidecar/metrics/main.go --help
```

**Note that their default values are all set to be compatible with Algorand's sandbox.**

#### RESTful API Service

Run the following command to start the API service:

```sh
make start-algosearch-backend

# OR this, which is the same command:
go run ./backend/app/algosearch/main.go --help
```

If you want to connect the API to sandbox, run:

```sh
make start-sandbox-algosearch-backend
```

#### Metric Service

Run the following command to start the metric service:

```sh
make start-algosearch-metrics

# OR this, which is the same command:
go run ./backend/app/sidecar/metrics/main.go
```

If you want to connect the metric to work with the sandbox, run:

```sh
make start-sandbox-algosearch-backend
```

#### Frontend



## Docker Support

To build AlgoSearch, run the following command and you will have an image `algosearch:1.1`:

```sh
make algosearch
```

You can also choose to build the containers separately:

```sh
# RESTful API
make algosearch-backend

# Metrics API
make algosearch-metrics

# Frontend
make algosearch-frontend
```


```shell
# kill local postgres to make way for sandbox's postgres
make kill-postgres

# delete couchdb data
sudo rm -rf db-data

# start couchdb
make start-local-couchdb

# insert the design documents into couchdb
make migrate-couch

# start sandbox

# start the service
make it-rain
```

# Heroku

https://devcenter.heroku.com/articles/container-registry-and-runtime#dockerfile-commands-and-runtime

docker tag algosearch-backend:latest kevguy/algosearch-backend:latest
docker push kevguy/algosearch-backend:latest
docker tag kevguy/algosearch-backend:latest registry.heroku.com/algosrch/web
docker push registry.heroku.com/algosrch/web
heroku container:release web --app algosrch
