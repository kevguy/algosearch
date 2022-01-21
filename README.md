# AlgoSearch

## Introduction

AlgoSearch is a open-sourced project that enables you to explore and search the Algorand blockchain for transactions, addresses, and blocks, assets, statistics, and more, in real-time. It's a simple, easy-to-deploy, and open-source block explorer to be used alongside an Algorand archival node.

It contains 3 services:

- Frontend app
  - The website of AlgoSearch
- RESTful API server
  - It connects to the Algorand archival node (and indexer, optional) and serves a set of API endpoints for the frontend to consume.
- Metrics server (optional)
  - It connects to the RESTful API and monitors its status

## Usage

### Prerequisities

Make sure you have a CouchDB database set up and a working Algorand node, ideally an archival node.

If you want to do tracing, you can set up Zipkin too.

#### CouchDB

If you want to start a CouchDB quickly in your local environment, you can run this command to start one using Docker:

```sh
make run-couch
```

- http://localhost:5984
  - username: `admin`
  - password: `password`
  - volume: `PROJECT_FOLDER/db-data`

### Getting Started

#### Using Docker

You can start the container with a [Docker image](https://hub.docker.com/r/kevguy/algosearch) which already contains all the three services (RESTful API, metrics and frontend):

```sh
docker run \
  -e ALGOSEARCH_WEB_ENABLE_SYNC=true \
  -e ALGOSEARCH_WEB_SYNC_INTERNAL=5s \
  -e ALGOSEARCH_COUCH_DB_HOST=234.567.89.0:5984 \
  -e ALGOSEARCH_COUCH_DB_USER=algorand \
  -e ALGOSEARCH_COUCH_DB_PASSWORD=algorand \
  -e ALGOSEARCH_COUCH_DB_NAME=algosearch \
  -e ALGOSEARCH_ALGOD_PROTOCOL=http \
  -e ALGOSEARCH_ALGOD_ADDR=234.567.89.0:4001 \
  -e ALGOSEARCH_ALGOD_TOKEN=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa \
  -e ALGOSEARCH_ZIPKIN_REPORTER_URI=http://234.567.89.0:9411/api/v2/spans \
  -e NEXT_PUBLIC_API_URL=http://0.0.0.0:5000 \
  -e METRICS_COLLECT_FROM=http://0.0.0.0:4000/debug/vars \
  algosearch:1.1
```

Please modify `NEXT_PUBLIC_API_URL` only when you are trying to connect to another backend.

Please modify `METRICS_COLLECT_FROM` only when you are trying to collect metrics from another RESTful API.

- Frontend: http://localhost:3000
- RESTful API: http://localhost:5000
- Metrics: http://localhost:3001

#### Using Docker-Compose

You can also use `docker-compose` to start all the services with each of them in separate Docker images.

Go inside `PROJECT_FOLDER/zarf/compose/compose-config.yaml` and change the environment variables accordingly, and then make use of these commands:

```sh
# Start everything using docker-compose
make up

# See the logs
make logs

# Stop the containers
make down
```

- Frontend: http://localhost:3000
- RESTful API: http://localhost:5000
- Metrics: http://localhost:3001
- Zipkin: http://localhost:9411
- CouchDB: http://localhost:5984

To build the docker images yourself:

```sh
# RESTful API
# algosearch-backend:1.1
make algosearch-backend
# algosearch-backend:latest
make algosearch-backend-latest

# Frontend
# algosearch-frontend:1.1
make algosearch-frontend
# algosearch-frontend:latest
make algosearch-frontend-latest

# Metrics
# algosearch-metrics:1.1
make algosearch-metrics
# algosearch-metrics:latest
make algosearch-metrics-latest
```

Additionally, here are some useful commands for Docker:

```sh
# Stop and remove all containers (not only AlgoSearch)
make docker-down-local

# See logging of all containers
make docker-logs-local

# Clean and remove all docker images
make docker-clean
```

### Local

To run AlgoSearch locally, you need to have the following dependencies:

- npm/yarn, for building and starting the frontend app
- golang, for building and starting the backend services
- a couchdb connection, for the backend RESTful API to store and retrieve data

### Installation

Install the dependencies for frontend and the other services:

```sh
# Install all the dependencies for RESTful API and metric services
make tidy

# Install dependencies for frontend app
cd frontend
yarn install
```

### CouchDB

If you haven't set up a database on CouchDB for AlgoSearch to use, run this command with the appropriate credentials to set it up:

```sh
go run backend/app/algo-admin/main.go \
		--couch-db-protocol=http \
		--couch-db-user=admin \
		--couch-db-password=password \
		--couch-db-host=0.0.0.0:5984 \
		--couch-db-name=algosearch \
		migrate
```

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

Start the API service:

```sh
go run ./backend/app/algosearch/main.go

# OR this, which is the same command but with
# better logging format
make start-algosearch-backend
```

If you are connecting to the API to sandbox, run:

```sh
make start-sandbox-algosearch-backend
```

#### Metric Service

Start the metric service:

```sh
go run ./backend/app/sidecar/metrics/main.go

# OR this, which is the same command but with
# better logging format
make start-algosearch-metrics
```

If you are connecting the metric service to work with the sandbox, run:

```sh
make start-sandbox-algosearch-backend
```

#### Frontend

Go inside the `frontend` folder:

```sh
cd frontend
yarn dev

# OR
yarn start
```

# Heroku

https://devcenter.heroku.com/articles/container-registry-and-runtime#dockerfile-commands-and-runtime

docker tag algosearch-backend:latest kevguy/algosearch-backend:latest
docker push kevguy/algosearch-backend:latest
docker tag kevguy/algosearch-backend:latest registry.heroku.com/algosrch/web
docker push registry.heroku.com/algosrch/web
heroku container:release web --app algosrch
