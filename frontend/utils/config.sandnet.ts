import { Config } from "./config";

const sandnet_config: Config = {
  "algod":{
    "host": "http://localhost",
    "port":4001,
    "token":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  },
  "indexer":{
    "host": "http://localhost",
    "port": 8980,
    "token":""
  },
  "network":"sandnet-v1"
}

export default sandnet_config;