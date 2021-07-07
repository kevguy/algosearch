/*
	initSync.js
	____________
	Provides scripts to test for and create blocks and transactions db.
*/

var constants = require('../global'); // Requiring global constants
var nano = require("nano")(`http://${constants.dbuser}:${constants.dbpass}@${constants.dbhost}`); // Setting up connection to db

/*
	Check for blocks database
	If it exists, continue. Else, create it.
*/
const initBlocksDB = () => {
	console.log("Checking for blocks database");

	return nano.db.list().then(body => {
		if (!body.includes('blocks')) {
			console.log("blocks database does not exist, creating...");
			return nano.db.create('blocks').then(insertLatestView);
		} else {
			console.log("blocks database exists, continuing...");
			return insertLatestView();
		}
	}).catch(error => {
		console.log("Exception when initializing blocks DB: " + error);
	})
};

/*
	Check for transactions database
	If it exists, continue. Else, create it.
*/
const initTransactionsDB = () => {
	console.log("Checking for transactions database");

	return nano.db.list().then(body => {
		if (!body.includes('transactions')) {
			console.log("transactions database does not exist, creating...");
			return nano.db.create('transactions').then(insertQueryView);
		} else {
			console.log("transactions database exists, continuing...");
			return insertQueryView();
		}
	}).catch(error => {
		console.log("Exception when initializing transactions DB: " + error);
	})
}

/*
	Check for addresses database
	If it exists, continue. Else, create it.
*/
const initAddressesDB = () => {
	console.log("Checking for addresses database");

	return nano.db.list().then(body => {
		if (!body.includes('addresses')) {
			console.log("addresses database does not exist, creating...");
			nano.db.create('addresses');
		} else {
			console.log("addresses database exists, continuing...");
		}
	}).catch(error => {
		console.log("Exception when initializing addresses DB: " + error);
	})
}

/*
	Insert the latest view into the blocks database.
	If it exists, continue. Else, create it.
 */
const insertLatestView = () => {
	console.log("Inserting Latest View");

	const db = nano.use("blocks");
	db.view("latest", "latest").then(body => {
		console.log("view exists, continuing...");
	}).catch(error => {
		db.insert({
			"views": {
				"latest": {
					"map": function(doc) { emit(doc.round); },
				},
			},
		}, "_design/latest").catch(error => {
			console.log("Exception when inserting latest view: " + error);
		});
	});
}

/*
	Insert the query view into the transactions database.
	If it exists, continue. Else, create it.
 */
const insertQueryView = () => {
	console.log("Inserting Query View");

	const db = nano.use("transactions");
	db.view("query", "bytimestamp").then(body => {
		console.log("view exists, continuing...");
	}).catch(error => {
		db.insert({
			"views": {
				"bytimestamp": {
					"map": function(doc) { emit(doc.round); },
				},
			},
		}, "_design/query").catch(error => {
			console.log("Exception when inserting query view: " + error);
		});
	});
}

async function init() {
	// Executing this file will also run the functions:
	await initBlocksDB();
	await initTransactionsDB();
	await initAddressesDB();
}

init();

// Export functions
module.exports = {
	initBlocksDB,
	initTransactionsDB,
	initAddressesDB,
	insertLatestView,
	insertQueryView,
};
