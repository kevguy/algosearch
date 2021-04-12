import { Config } from "./config";

const testnet_config: Config = {
  "algod":{
    "host": "https://testnet.algoexplorerapi.io",
    "port": 0,
    "token":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  },
  "indexer":{
    "host": "https://testnet.algoexplorerapi.io/idx2",
    "port": 0,
    "token":""
  },
  "network":"TestNet"
}

export default testnet_config;