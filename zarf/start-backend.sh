#!/bin/sh

if [ "$ALGOSEARCH_COUCH_DB_INIT" = "true" ];
then
  echo "Initializing Couch DB Database..."
  # https://github.com/apache/couchdb-docker/issues/55
  curl -X PUT "${ALGOSEARCH_COUCH_DB_PROTOCOL}://${ALGOSEARCH_COUCH_DB_USER}:${ALGOSEARCH_COUCH_DB_PASSWORD}@${ALGOSEARCH_COUCH_DB_HOST}/_users" && \
  curl -X PUT "${ALGOSEARCH_COUCH_DB_PROTOCOL}://${ALGOSEARCH_COUCH_DB_USER}:${ALGOSEARCH_COUCH_DB_PASSWORD}@${ALGOSEARCH_COUCH_DB_HOST}/_users/_security" -d '{}' && \
  curl -X PUT "${ALGOSEARCH_COUCH_DB_PROTOCOL}://${ALGOSEARCH_COUCH_DB_USER}:${ALGOSEARCH_COUCH_DB_PASSWORD}@${ALGOSEARCH_COUCH_DB_HOST}/algo_global?partitioned=false"
##      curl -X PUT http://algorand:algorand@algosearch-db:5984/algo_global?partitioned=false &&
##      curl -X PUT http://algorand:algorand@algosearch-db:5984/algo_global/_security -d '{\"members\": {}, \"admins\": {\"roles\": [\"_admin\"] }}'"
  sleep 5s
fi

if [ "$ALGOSEARCH_COUCH_DB_MIGRATE" = "true" ];
then
  echo "Migrating Couch DB Database..."
  ./admin migrate
  sleep 5s
fi

if [ "$ALGOSEARCH_BACKEND_DISABLED" != "true" ]; then ./algosearch; fi & \
