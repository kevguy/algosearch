const fetch = require('isomorphic-fetch')


// const url = 'http://kevin:makechesterproud!@89.39.110.254:5984/algo_global/_design/block/_view/blockByRoundNoCount?inclusive_end=true&start_key=400000&end_key=300000&descending=true&reduce=true&group_level=1&skip=0&limit=101'

const url = 'http://kevin:makechesterproud!@89.39.110.254:5984/algo_global/_design/block/_view/blockByRoundNoCount?inclusive_end=true&start_key=3100000&end_key=3200000&reduce=true&group_level=1&skip=0'


for (let i = 2000000; i < 3000000; i+=10000) {
    const url1 = `http://kevin:makechesterproud!@89.39.110.254:5984/algo_global/_design/block/_view/blockByRoundNoCount?inclusive_end=true&start_key=${i}&end_key=${i+10000}&reduce=true&group_level=1&skip=0`
    fetch(url1)
        .then(res => res.json())
        .then(res => {
	    // console.log(res);
            console.log(`Key from ${i} to ${i+10000}: ${res.rows.length}`);
        });
}

// fetch(url)
//     .then(res => res.json())
//     .then(res => {
//         console.log(res.rows.length);
//     });
// 
