/*
	Query: addresses.js
	____________
	Addresses endpoint
	____________
	Various Return schemas
*/

var constants = require('../global'); // Require global constants
const axios = require("axios"); // Axios for requests

// Export express routes
module.exports = function(app) {

	// --> Single address data retrieval
	app.get('/addressservice/:address', function(req, res) {
		const address = req.params.address; // Get address from url

		// Request basic account information
		axios({
			method: 'get',
			url: `${constants.algodurl}/v2/accounts/${address}`, // Request transaction details endpoint
			headers: {'X-Algo-API-Token': constants.algodapi}
		}).then(response => {
			let result = response.data; // Set data to result

			axios({
				method: 'get',
				url: `${constants.algoIndexerUrl}/v2/accounts/${address}/transactions?limit=25`,
				headers: {'X-Indexer-API-Token': constants.algoIndexerToken}
			}).then(async resp => {
				let rounds = [];

				for (let i = 0; i < resp.data.transactions.length; i++) {
					if (resp.data.transactions[i]['confirmed-round']) {
						rounds.push(resp.data.transactions[i]['confirmed-round']);
					}
				}

				let uniqueRounds = [...new Set(rounds)];
				let roundsWithTimestamp = [];

				for (let i = 0; i < uniqueRounds.length; i++) {
					await axios({
						method: 'get',
						url: `${constants.algoIndexerUrl}/v2/blocks/${uniqueRounds[i]}`,
						headers: {'X-Indexer-API-Token': constants.algoIndexerToken}
					}).then(blockresponse => {
						roundsWithTimestamp.push({"round": blockresponse.data.round, "timestamp": blockresponse.data.timestamp});
					}).catch(error => {
						res.status(501);
						console.log("Exception when querying for round timestamp: " + error);
					})
				}

				for (let i = 0; i < resp.data.transactions.length; i++) {
					let round = resp.data.transactions[i]['confirmed-round'];
					for (let j = 0; j < roundsWithTimestamp.length; j++) {
						if (roundsWithTimestamp[j].round === round) {
							resp.data.transactions[i].timestamp = roundsWithTimestamp[j].timestamp;
						}
					}
				}

				// Add transactions to result
				result.confirmed_transactions = resp.data.transactions;
				res.send(result);
			}).catch(error => {
				res.status(501);
				console.log("Exception when retrieving address transactions: " + error);
			});
		}).catch(error => {
			res.status(501);
			console.log("Exception when retrieving address information: " + error);
		});
	});
}
