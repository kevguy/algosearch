import { TxType } from "../utils/stringUtils";
export interface ICurrentRoundResponse {
  round: number;
  "genesis-id": number;
  transactions?: [];
}

export interface IAsaResponse {
  index: number;
  params: {
    clawback: string;
    creator: string;
    decimals: number;
    freeze: string;
    manager: string;
    name: string;
    reserve: string;
    total: number;
    "unit-name": string;
    url: string;
  }
}

export type TransactionResponse = {
  id: number;
  "genesis-id": number;
  "genesis-hash": string;
  "confirmed-round": number;
  "tx-type": TxType;
  sender: string;
  "sender-rewards": number;
  "receiver-rewards": number;
  "payment-transaction": {
    amount: number;
    "close-amount": number;
    "close=remainder-to": string;
    receiver: string;
  };
  "asset-transfer-transaction": {
    "asset-id": number;
    amount: number;
    receiver: string;
    "close-to": string;
    sender: string;
  };
  "asset-config-transaction": {
    params: {
      creator: string;
      decimals: number;
      total: number;
      // more info if type is Asset Config
      manager?: string;
      reserve?:string;
      freeze?: string;
      clawback?: string;
      "metadata-hash"?: string;
      name?:string;
      total?: number;
      "unit-name"?: string;
      url?: string;
    };
  };
  fee: number;
  "round-time": number;
  "first-valid": number;
  "last-valid": number;
  timestamp: number;
  note: string;
};

export interface IBlockRewards {
  "fee-sink": string;
  "rewards-calculation-round": number;
  "rewards-level": number;
  "rewards-pool": string;
  "rewards-rate": number;
  "rewards-residue": number;
}

export interface IBlockResponse {
  "block-hash": string;
  doc_type: string;
  "genesis-hash": string;
  "genesis-id": string;
  "previous-block-hash": string;
  proposer: string;
  rewards: IBlockRewards;
  round: number;
  seed: string;
  timestamp: number;
  transactions?: TransactionResponse[];
  "transaction-root": string;
  "txn-counter": number;
  "_upgrade-state": {
    "current-protocol": string;
  };
  "upgrade-vote": {};
  _id: string;
  _rev: string;
}

export interface ILatestBlocksResponse {
  num_of_blks: number;
  num_of_pages: number;
  items: [];
}

export interface ISupply {
  current_round: number;
  "online-money": string;
}
