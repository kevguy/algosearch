/*
	syncAll.js
	____________
	Syncs current round with blocks and transactions database.
*/
const algosdk = require('algosdk');
var constants = require('./service/global'); // Require global constants
const axios = require('axios'); // Require axios for requests
// const nano = require("nano")(`http://${constants.dbuser}:${constants.dbpass}@${constants.dbhost}`); // Connect nano to db

// let blocks = nano.db.use('blocks'); // Connect to blocks db
// let transactions = nano.db.use('transactions'); // Connect to transactions db
// let addresses = nano.db.use('addresses'); // Connect to addresses db

const algoUrl = new URL(constants.algodurl);
const client = new algosdk.Algodv2(
    constants.algodapi,
    algoUrl,
    algoUrl.port ? algoUrl.port : 8080);


async function getProposerAndBlockHash(client, blockNum) {
    try {
        const blk = await client.block(blockNum).do();
        const proposer = algosdk.encodeAddress(blk["cert"]["prop"]["oprop"]);
        const blockHash = Buffer.from(blk["cert"]["prop"]["dig"]).toString("base64");
        return { proposer, blockHash };
    } catch (e) {
        console.log("[getProposerAndBlockHash]: Error getting proposer: " + e);
    }
}
/*
	Get current round number from algod
*/
async function getCurrentRound() {
    let round;

    await axios({
        method: 'get',
        url: `${constants.algodurl}/v2/ledger/supply`, // Request /ledger/supply endpoint
        headers: {'X-Algo-API-Token': constants.algodapi}
    }).then(response => {
        round = response.data.current_round; // Collect round
    }).catch(error => {
        console.log("[getCurrentRound]: Exception when getting current round: " + error);
    })

    return round;
}

async function getThis() {
    const response = await axios({
        method: 'get',
        // url: `${constants.algoIndexerUrl}/v2/blocks/${blockNum + increment + 1}`, // Retrieve each block in succession
        url: `${constants.algoIndexerUrl}/v2/blocks/5300`, // Retrieve each block in succession
        headers: {'X-Indexer-API-Token': constants.algoIndexerToken}
    });

    // const { proposer, blockHash }= await getProposerAndBlockHash(client, blockNum + increment + 1);
    const { proposer, blockHash }= await getProposerAndBlockHash(client, 5300);
    console.log({
        ...response.data,
        proposer,
        blockHash,
    }); // Push block to array
}

getThis();
