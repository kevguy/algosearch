# AlgoSearch


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

## Folder Structure

```
- backend
- vendor
- zarf
- dockerfile.algosearch-backend.dockerignore
- dockerfile.metrics.dockerignore
- go.mod
- makefile
- README.md
```

- backend: contains all the Go application code that's required to run the backend api
- vendor: contains all the Go dependencies the project needs to run the backend
- zarf: contains all the Docker and Docker compose related files
- dockerfile.algosearch-backend.dockerignoe: the `.dockerignore` file needed for building the `algosearch-backend` image
- dockerfile.metrics.dockerignoe: the `.dockerignore` file needed for building the `metrics` image
- go.mod: contains the list of Go dependencies the backend needs
- makefile: contains a general list of useful commands to run the project

## Approach

Every piece of data is stored inside a global database in CouchDB, which include the following entities:

- block (which contains the transactions)
- transactions
- accounts
- applications
- assets

Thus we have the following views:

```js
const views = {
  "views": {
    "block": function(doc) {
      if (doc.type === 'block') {
        emit(doc.round, doc);
      }
    },
    "txn": function(doc) { 
      if (doc.type === 'txn') {
        emit(doc.id, doc);
      }
    },
    "acct": function(doc) {
      if (doc.type === 'acct') {
        emit([doc.id, 0], doc);
      } else if (doc.type === 'acct_txn') {
        emit([doc.id, 1], doc) 
      }
    },
    "asset": function(doc) {
      if (doc.type === 'asset') {
        emit([doc.id, 0], doc);
      } else if (doc.type === 'asset_txn') {
        emit([doc.id, 1], doc)
      }
    },
    "app": function(doc) {
      if (doc.type === 'app') {
        emit([doc.id, 0], doc);
      } else if (doc.type === 'app_txn') {
        emit([doc.id, 1], doc)
      } 
    }
  },
};
```

The prefixes may not be necessary, but there may be conflicts between different entity IDs

For every block retrieved:
    - add the type `block` in it and store it using the `block_round_number` as the id  (using `block_` as the prefix)
    - extract the transactions, and for each transaction inside:
        - store the transaction (add type `transaction` into the document) using `transaction_transaction_id` as the id
        - extract all the addresses we could find, which may be of account, application or asset
        - test each address with Indexer/Algod to find out if it's account, application or asset
        - using account as an example:
            - insert a document with the type `acct` using key `['acct_account_id', 0]`
            - insert all the associated transactions with the key `['acct_account_id', 1]`


There are 5 kinds of transactions:
- PaymentTx (PaymentTransaction)
- KeyRegistrationTx (KeyregTransaction)
- AssetConfigTx (AssetConfigTransaction)
- AssetTransferTx (AssetTransferTransaction)
- AssetFreezeTx (AssetFreezeTransaction)
- ApplicationCallTx (ApplicationTransaction)

For all the transactions, there's a sender

### Payment Transaction

account document --> ([account_addr, 0], null)
transaction document --> ([account_addr, 1], null)



# Introduction to Views

Views are useful for many purposes:
- Filtering the documents in your database to find those relevant to a particular process.
- Extracting data from your documents and presenting it in a specific order.
- Building efficient indexes to find documents by any value or structure that resides in them.
- Use these indexes to represent relationships among documents.
- Finally, with views you can make all sorts of calculations on the data in your documents. For example, if documents represent your company's financial transactions, a view can answer the question of what the spending was in the last week, month or year.


## What is a View?

Let;s go through the different use cases. First is extracting data that you might need for a special purpose in a specific order. For a front page, we want a list of blog post titles sorted by date. We'll work with a set of example documents as we walk through how views work:

```json
{
    "_id":"biking",
    "_rev":"AE19EBC7654",

    "title":"Biking",
    "body":"My biggest hobby is mountainbiking. The other day...",
    "date":"2009/01/30 18:04:11"
}
```

```json
{
    "_id":"bought-a-cat",
    "_rev":"4A3BBEE711",

    "title":"Bought a Cat",
    "body":"I went to the the pet store earlier and brought home a little kitty...",
    "date":"2009/02/17 21:13:39"
}
```

```json
{
    "_id":"hello-world",
    "_rev":"43FBA4E7AB",

    "title":"Hello World",
    "body":"Well hello and welcome to my new blog...",
    "date":"2009/01/15 15:52:20"
}
```

Three will do for the example. Note that the documents are sorted by “_id”, which is how they are stored in the database. Now we define a view. Bear with us without an explanation while we show you some code:

```js
function(doc) {
    if (doc.date && doc.title) {
        emit(doc.date, doc.title);
    }
}
```

This is a **map** function, and it is written in JavaScript. If you are not familiar with JavaScript but have used C or any other C-like language such as Java, PHP, or C#, this should look familiar. It is a simple function definition.

You provide CouchDB with view functions as strings stored inside the `views` field of a design document. You don’t run it yourself. Instead, when you query your view, CouchDB takes the source code and runs it for you on every document in the database your view was defined in. You query your view to retrieve the **view** result.

The `emit()` function always takes two arguments: the first is `key`, and the second is `value`. The `emit(key, value)` function creates an entry in our **view result**. One more thing: the `emit()` function can be called multiple times in the map function to create multiple entries in the view results from a single document, but we are not doing that yet.

CouchDB takes whatever you pass into the `emit()` function and puts it into a list (see Table 1, "View Results" below). Each row in that list includes the **key** and **value**. More importantly, the list is sorted by key (by `doc.date` in our case). The most important feature of a view result is that it is sorted by **key**. We will come back to that over and over again to do neat things. Stay tuned.

Table 1. View results:

| Key | Value |
|:---:|:---:|
| “2009/01/15 15:52:20”	| “Hello World” |
| “2009/01/30 18:04:11”	| “Biking” |
| “2009/02/17 21:13:39”	| “Bought a Cat” |

When you query your view, CouchDB takes the source code and runs it for you on every document in the database. If you have a lot of documents, that takes quite a bit of time and you might wonder if it is not horribly inefficient to do this. Yes, it would be, but CouchDB is designed to avoid any extra costs: it only runs through all documents once, when you first query your view. If a document is changed, the map function is only run once, to recompute the keys and values for that single document.

The view result is stored in a B-tree, just like the structure that is responsible for holding your documents. View B-trees are stored in their own file, so that for high-performance CouchDB usage, you can keep views on their own disk. The B-tree provides very fast lookups of rows by key, as well as efficient streaming of rows in a key range. In our example, a single view can answer all questions that involve time: "Give me all the blog posts from last week" or "last month" or "this year". Pretty neat.

When we query our view, we get back a list of all documents sorted by date. Each row also includes the post title so we can construct links to posts. Table 1 is just a graphical representation of the view result. The actual result is JSON-encoded and contains a little more metadata:

```json
{
    "total_rows": 3,
    "offset": 0,
    "rows": [
        {
            "key": "2009/01/15 15:52:20",
            "id": "hello-world",
            "value": "Hello World"
        },

        {
            "key": "2009/01/30 18:04:11",
            "id": "biking",
            "value": "Biking"
        },

        {
            "key": "2009/02/17 21:13:39",
            "id": "bought-a-cat",
            "value": "Bought a Cat"
        }

    ]
}
```



Rmb to support Group Transactions, this seems hard tbh.


References:
- https://developer.algorand.org/docs/reference/node/config/
- https://developer.algorand.org/tutorials/betanet-sandbox/#1-running-betanet-for-the-first-time


TODO:
- transactions by account
- sync blocks
- ASA names
- add original swagger doc
- get algod ledger supply (https://developer.algorand.org/docs/rest-apis/algod/v2/#get-v2ledgersupply)
  - latest/current round
  - online money
  - circulating supply
