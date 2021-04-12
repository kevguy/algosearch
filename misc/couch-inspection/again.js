const fetch = require('isomorphic-fetch');
fetch("https://cal-engine.choco-up.com/v1/bank_rev_settlements/74493197-307c-43bb-a405-08506982ba46/appeal", {
  "headers": {
    "accept": "*/*",
    "accept-language": "en-GB,en-US;q=0.9,en;q=0.8",
    "sec-ch-ua": "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"",
    "sec-ch-ua-mobile": "?0",
    "sec-ch-ua-platform": "\"macOS\"",
    "sec-fetch-dest": "empty",
    "sec-fetch-mode": "cors",
    "sec-fetch-site": "same-site"
  },
  "referrer": "https://app.choco-up.com/",
  "referrerPolicy": "strict-origin-when-cross-origin",
  "body": null,
  "method": "PUT",
  "mode": "cors",
  "credentials": "omit"
}).then(res => res.json())
.then(res => {
	console.log(res);
});
