export interface ICurrentRoundResponse {
  round: number;
  "genesis-id": number;
  transactions?: [];
}

export interface IBlockResponse {
  "block-hash": string;
  doc_type: string;
  "genesis-hash": string;
  "genesis-id": string;
  "previous-block-hash": string;
  proposer: string;
  rewards: {
    "fee-sink": string;
    "rewards-calculation-round": number;
    "rewards-level": number;
    "rewards-pool": string;
    "rewards-rate": number;
    "rewards-residue": number;
  };
  round: number;
  seed: string;
  timestamp: number;
  transactions?: [];
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
