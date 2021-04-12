<p>
<img src="https://i.imgur.com/dsBUUav.png" width="300" alt="AlgoSearch logo image" />
</p>

[![License](https://img.shields.io/badge/License-Apache%202.0-yellowgreen.svg)](https://opensource.org/licenses/Apache-2.0)

# AlgoSearch ([live deployment](https://algosearch.io))
AlgoSearch enables you to explore and search the [Algorand blockchain](https://www.algorand.com/) for transactions, addresses, blocks, assets, statistics, and more, in real-time. It's a simple, easy-to-deploy, and open-source block explorer to be used alongside an Algorand archival node.

**Dependencies**
* [Node.js](https://nodejs.org/en/) 8+ for use with server and front-end.
* [go-algorand](https://github.com/algorand/go-algorand) for Algorand `goal` node (must support archival indexing).
* [Algorand Indexer](https://github.com/algorand/indexer) for reading committed blocks from the Algorand blockchain and maintains a database of transactions and accounts that are searchable and indexed.
* [CouchDB](https://couchdb.apache.org/) as database solution.

Work on AlgoSearch is funded by the [Algorand Foundation](https://algorand.foundation) through a grant to [Anish Agnihotri](https://github.com/anish-agnihotri). The scope of work includes the development of an open-source block explorer (AlgoSearch) and a WIP analytics platform.

# Run locally

## Linux / OSX
The [go-algorand](https://github.com/algorand/go-algorand) node currently aims to support only Linux and OSX environments for development.

## Disclaimer
Simpler installation instructions, a hands-on guide, and a one-click deploy Docker image will be published upon final completion of AlgoSearch.

## Step 1: Environment setup

This section explains how to set up everything locally.

### The Native Approach

#### Algorand's Node

First you'll need to install [Algorand's Node](https://developer.algorand.org/docs/run-a-node/setup/install/) locally. Follow the instructions through the hyperlink.

Make sure node is running on the preferred network and that algod details are correct in `service/global.js`.

#### Indexer

Then you'll need to install the [Indexer](https://developer.algorand.org/docs/run-a-node/setup/indexer/) locally. Follow the instructions through the hyperlink.

#### CouchDB

Finally you'll need to install [CouchDB](https://docs.couchdb.org/en/stable/install/index.html) locally.

You can also run CouchDB using Docker easily:

```sh
# Create a folder called db-data
mkdir db-data

docker run -e COUCHDB_USER=admin -e COUCHDB_PASSWORD=password -p 5984:5984 --name my-couchdb -v $(pwd)/db-data:/opt/couchdb/data -d couchdb
```

If you are using docker compose to start the services, you can skip this step.

### The Sandbox Approach

You can also set up **Algorand's Node** and the **Indexer** using [Algorand's Sandbox](https://github.com/algorand/sandbox). Follow the instructions through the hyperlink.

Note that you will still have to set up CouchDB if you are not using the `docker-compose.yml` offered here. 

## Step 2: Configuration

1. Enter your site name in `src/constants.js`.
    - it's set to be `http://localhost:8000` as default, but if you are changing the port, remember to update `PORT` in `server.js` and `server.local.js`
2. Enter the API endpoint of the Algorand's Node in `src/constants.js`.
3. Copy `service/global.sample.js` to `service/global.js` and enter your node and DB details.
    - If you are using the Sandbox approach, copy `service/global.sandbox.js` instead and update the CouchDB details if needed.

## Step 3: Running AlgoSearch

### The Native Approach

#### Install the dependencies

```
# Run in folder root directory
npm install
```

#### Build the code

```
# Run in folder root directory
npm run build
```

#### Run it

Make sure the configurations in your `src/global.js` is correct, then you'll have to do three things.

One, execute the following to create tables in CouchDB:

```sh
node service/sync/initSync.js
```

Second, execute the following to start syncing the tables:

```sh
node service/sync/syncAll.js
```

Note that this step takes time to sync and should stay running as long as the server is running.

Finally, start the server:

```sh
nodemon server.js
```

### Docker Approach

You can skip the native approach entirely and simply start the application using Docker (remember to make sure your `src/global.js` is having the correct details):

```sh
# Build the image
docker build -t algosearch .

# Run the container
docker run algosearch
```

If you are using Linux and your container needs to access the host machine, for instance, the CouchDB you set up on your machine, run the following the start the container:

```
docker run --network="host" algosearch
```

#### Docker Compose Approach

To start the server using `docker-compose`, you only need the Node and Indexer, and use DB details in `src/global.sandbox.js`, that is, make sure

```
dbhost = 'couchdb.server:5984', // Database URL
dbuser = 'admin', // Database user
dbpass = 'password', // Database password
```

Then start the services:

```sh
# Create the folder for CouchDB
mkdir db-data

# Start the services
docker-compose up
```

If your Node and Indexer are on the host machine, your containers will have to access `localhost`, instead of using `docker-compose up`, run the startup script instead:

```sh
bash docker-run.sh
```

This script will find the IPs of `localhost` and to be accessed through `dockerhost`. In other words, in your `src/global.js`, use `dockerhost` instead of `localhost`.

# The Fastest Approach

The fastest approach to set everything up for development is using the Sandbox and docker-compose. To do that, just setup the Sandbox and do the following:

```
# Create the folder for CouchDB
mkdir db-data

# Start the services
bash docker-run.sh
```

# Documentation
The [Wiki](https://github.com/Anish-Agnihotri/algosearch/wiki) is currently under construction.

# License
[![License](https://img.shields.io/badge/License-Apache%202.0-yellowgreen.svg)](https://opensource.org/licenses/Apache-2.0)

Copyright (C) 2020, Anish Agnihotri.
