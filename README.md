# AlgoSearch

## Introduction

AlgoSearch is an open-sourced project that enables you to explore and search the Algorand blockchain for transactions, blocks, addresses, assets, statistics, and more, in real-time. It's a simple, easy-to-deploy, and open-source block explorer to be used alongside an Algorand archival node.

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
  -e NEXT_PUBLIC_ALGOD_PROTOCOL=http \
  -e NEXT_PUBLIC_ALGOD_ADDR=0.0.0.0:4001 \
  -e NEXT_PUBLIC_ALGOD_TOKEN=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa \
  -e METRICS_COLLECT_FROM=http://0.0.0.0:4000/debug/vars \
  algosearch:1.1
```

Please modify `NEXT_PUBLIC_API_URL` only when you are trying to connect to another backend.

Please modify `METRICS_COLLECT_FROM` only when you are trying to collect metrics from another RESTful API.

`NEXT_PUBLIC_ALGOD_PROTOCOL`, `NEXT_PUBLIC_ALGOD_ADDR`, and `NEXT_PUBLIC_ALGOD_TOKEN` are needed for disassembly of LogicSig, approval program, and clear state program on the transaction page. The feature is only available when `NEXT_PUBLIC_ALGOD_ADDR` contains `0.0.0.0` or `127.0.0.1` or `localhost`.

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

### Running Locally

To run AlgoSearch locally, you need to have the following dependencies:

- npm/yarn, for building and starting the frontend app
- golang, for building and starting the backend services
- a couchdb connection, for the backend RESTful API to store and retrieve data

#### Installation

Install the dependencies for frontend and the other services:

```sh
# Install all the dependencies for RESTful API and metric services
make tidy

# Install dependencies for frontend app
cd frontend
yarn install
```

#### CouchDB

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

#### Backend

Both the restful API and metric services are configurable. Run the following commands to see what variables that can be configured through command line arguments or environment variables:

```sh
# RESTful API
go run ./backend/app/algosearch/main.go --help

# Metrics
go run ./backend/app/sidecar/metrics/main.go --help
```

**Note that their default values are all set to be compatible with Algorand's sandbox.**

##### RESTful API Service

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

##### Metric Service

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
yarn build
yarn start
```

## Core Team

<table>
  <tbody>
    <tr>
      <td align="center" width="33.3%" valign="top">
        <img width="150" height="150" src="https://github.com/kevguy.png?s=150">
        <br>
        <a href="https://github.com/kevguy">Kevin Lai</a>
        <p>Core Services</p>
        <br>
        <p>Golang, Linkin Park, South Park, and Red Bull</p>
      </td>
      <td align="center" width="33.3%" valign="top">
        <img width="150" height="150" src="https://media.giphy.com/media/OXY1YM1QSKt23KAg1x/giphy.gif">
        <br>
        <a href="https://github.com/fionnachan">Fionna Chan</a>
        <p>Frontend & UI/UX Design</p>
        <br>
        <p>Making the world a better place with OSS, one line at a time</p>
      </td>
      <td align="center" width="33.3%" valign="top">
        <img width="150" height="150" src="https://github.com/Uppers.png?s=150">
        <br>
        <a href="https://github.com/Uppers">Thomas Upfield</a>
        <p>Documentation & Business Relations</p>
        <br>
        <p>Algorand Evangelist. DeFi, tokenomics, and analytics</p>
      </td>
     </tr>
  </tbody>
</table>

## Special Thanks to

- [@ardanlabs](https://github.com/ardanlabs) for [service](https://github.com/ardanlabs/service), which taught us everything we know about Golang and offering a well-designed sample API service as our foundation.

- [@Anish-Agnihotri](https://github.com/Anish-Agnihotri) for his contribution to the original [AlgoSearch](https://github.com/Anish-Agnihotri/algosearch) written with create-react-app and a Node.js backend.

## Licensing

```
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
