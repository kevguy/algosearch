import Configuration from "./config.testnet";

export interface ApiConf {
    host: string
    port: number | string
    token: string
}

export interface Config {
    algod:  ApiConf 
    indexer:  ApiConf
    network: string
}

export default Configuration
