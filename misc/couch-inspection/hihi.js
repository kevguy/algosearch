const fetch = require('isomorphic-fetch')

// fetch("https://cal-engine-uat.choco-up.com/v1/remittance_streams/set_channel_accts", {
//   "headers": {
//     "accept": "*/*",
//     "accept-language": "en-GB,en-US;q=0.9,en;q=0.8",
//     "content-type": "text/plain;charset=UTF-8",
//     "sec-ch-ua": "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"",
//     "sec-ch-ua-mobile": "?0",
//     "sec-ch-ua-platform": "\"macOS\"",
//     "sec-fetch-dest": "empty",
//     "sec-fetch-mode": "cors",
//     "sec-fetch-site": "cross-site"
//   },
//   "referrer": "http://localhost:3000/",
//   "referrerPolicy": "strict-origin-when-cross-origin",
//   "body": "{\"stream_id\":\"fd12c9d8-d5c4-468f-bc3d-fd88d2c290a1\",\"channel\":\"bank\",\"accounts\":{\"xero\":[\"4ff3856b-5427-4e37-afdf-bfdb11a13a12\"],\"zohobooks\":[],\"quickbooks\":[\"Turnkey Lender  , Turnkey Lender Pte Ltd\"],\"otherBankAccounts\":[]}}",
//   "method": "PUT",
//   "mode": "cors",
//   "credentials": "omit"
// });

// fetch("https://cal-engine-uat.choco-up.com/v1/remittance_streams/set_channel_accts", {
//   "headers": {
//     "accept": "*/*",
//     "accept-language": "en-GB,en-US;q=0.9,en;q=0.8",
//     "content-type": "text/plain;charset=UTF-8",
//     "sec-ch-ua": "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"",
//     "sec-ch-ua-mobile": "?0",
//     "sec-ch-ua-platform": "\"macOS\"",
//     "sec-fetch-dest": "empty",
//     "sec-fetch-mode": "cors",
//     "sec-fetch-site": "cross-site"
//   },
//   "referrer": "http://localhost:3000/",
//   "referrerPolicy": "strict-origin-when-cross-origin",
//   "body": "{\"stream_id\":\"fd12c9d8-d5c4-468f-bc3d-fd88d2c290a1\",\"channel\":\"payment\",\"accounts\":{\"stripe\":[\"Grain\"],\"shopify\":[],\"paypal\":[\"abc\"],\"bbmsl\":[]}}",
//   "method": "PUT",
//   "mode": "cors",
//   "credentials": "omit"
// });

// fetch("https://cal-engine-uat.choco-up.com/v1/remittance_streams/set_channel_accts", {
//   "headers": {
//     "accept": "*/*",
//     "accept-language": "en-GB,en-US;q=0.9,en;q=0.8",
//     "content-type": "text/plain;charset=UTF-8",
//     "sec-ch-ua": "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"",
//     "sec-ch-ua-mobile": "?0",
//     "sec-ch-ua-platform": "\"macOS\"",
//     "sec-fetch-dest": "empty",
//     "sec-fetch-mode": "cors",
//     "sec-fetch-site": "cross-site"
//   },
//   "referrer": "http://localhost:3000/",
//   "referrerPolicy": "strict-origin-when-cross-origin",
//   "body": "{\"stream_id\":\"fd12c9d8-d5c4-468f-bc3d-fd88d2c290a1\",\"channel\":\"analytic\",\"accounts\":{}}",
//   "method": "PUT",
//   "mode": "cors",
//   "credentials": "omit"
// });

fetch("https://cal-engine-uat.choco-up.com/v1/remittance_streams/set_channel_accts", {
  "headers": {
    "accept": "*/*",
    "accept-language": "en-GB,en-US;q=0.9,en;q=0.8",
    "content-type": "text/plain;charset=UTF-8",
    "sec-ch-ua": "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"",
    "sec-ch-ua-mobile": "?0",
    "sec-ch-ua-platform": "\"macOS\"",
    "sec-fetch-dest": "empty",
    "sec-fetch-mode": "cors",
    "sec-fetch-site": "cross-site"
  },
  "referrer": "http://localhost:3000/",
  "referrerPolicy": "strict-origin-when-cross-origin",
  "body": "{\"stream_id\":\"fd12c9d8-d5c4-468f-bc3d-fd88d2c290a1\",\"channel\":\"sales\",\"accounts\":{\"airwallex\":[],\"app-store\":[],\"ebay\":[],\"google-play-store\":[],\"lazada\":[],\"shopify\":[\"shekou-woman\",\"shekou-festive\",\"shekou-woman-nz\"],\"woocommerce\":[\"https://car-learning.com\"]}}",
  "method": "PUT",
  "mode": "cors",
  "credentials": "omit"
});
