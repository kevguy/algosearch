# AlgoSearch

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

