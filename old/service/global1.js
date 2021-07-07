// Global exports
module.exports = Object.freeze({
    // Supply your own CouchDB configs
    // No need to change the following if you are using the default docker-compose.yml
    dbhost: 'couchdb.server:5984', // Database URL
    dbuser: 'admin', // Database user
    dbpass: 'password', // Database password
    // No need to change the following
    algodurl: 'http://localhost:4001', // Algod node endpoint
    algodapi: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa', // Algod node API access header
    algoIndexerUrl: 'http://localhost:8980', // Algo Indexer endpoint
    algoIndexerToken: '',
});
