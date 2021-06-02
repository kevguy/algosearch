var constants = require('../global'); // Require global constants
const axios = require('axios'); // Axios for requests

module.exports = function(app) {

	app.get('/detect/:string', function(req, res) {
		const string = req.params.string;

		axios({
			method: 'get',
			url: `${constants.algodurl}/v2/accounts/${string}`,
			headers: {'X-Algo-API-Token': constants.algodapi}
		}).then(addrresp => {
			if (addrresp.data.address === string) {
				res.send('address');
			} else {
				res.send('error');
				console.log("Exception: false address formatting");
			}
		}).catch((e) => {
			axios({
				method: 'get',
				url: `${constants.algoIndexerUrl}/v2/blocks/${string}`,
				headers: {'X-Indexer-API-Token': constants.algoIndexerToken}
			}).then(blockresp => {
				if (blockresp.data.round.toString() === string) {
					res.send('block');
				} else {
					res.send('error');
					console.log("Exception: false block formatting");
				}
			}).catch(() => {
				axios({
					method: 'get',
					url: `${constants.algoIndexerUrl}/v2/transactions/${string}`,
					headers: {'X-Indexer-API-Token': constants.algoIndexerToken}
				}).then(transresp => {
					if (transresp.data.transaction && transresp.data.transaction.id === string) {
						res.send('transaction');
					} else {
						res.send('error');
						console.log("Exception: false transaction formatting");
					}
				}).catch(error => {
					res.send('error');
					console.log("Exception: cannot detect parsed string: ", error);
				})
			})
		})
	});
}
