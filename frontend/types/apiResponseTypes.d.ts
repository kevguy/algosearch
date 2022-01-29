import { TxType } from "../utils/stringUtils";

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

type multisigSubsig = {
  "public-key": string;
  signature?: string;
}

export type AccountOwnedAsset = {
  amount: number;
  "asset-id": number;
  creator: string;
  "is-frozen": boolean;
}

export type AccountResponse = {
  address: string;
  amount: number;
  "amount-without-pending-rewards": number;
  "apps-total-schema": { "num-byte-slice": number; "num-uint": number;},
  assets: AccountOwnedAsset[];
  "created-assets":IAsaResponse[];
  participation: {
    "selection-participation-key": string | null;
    "vote-first-valid": number;
    "vote-key-dilution": number;
    "vote-last-valid": number;
    "vote-participation-key": string | null;
  };
  "pending-rewards": number;
  "reward-base": number;
  "rewards": number;
  "round": number;
  "status": "Online" | "Offline";
}


export type AccountTxsResponse = {
  "num_of_pages": number,
  "num_of_txns": number,
  items: TransactionResponse[]
}

export type TransactionResponse = {
  id: string;
  group?: string;
  "genesis-id": string;
  "genesis-hash": string;
  "confirmed-round": number;
  "tx-type": TxType;
  sender: string;
  "sender-rewards": number;
  "receiver-rewards": number;
  "application-transaction": {
    accounts: string[];
    "application-args": [];
    "application-id": number;
    "approval-program"?: string;
    "foreign-apps": [];
    "foreign-assets": number[];
    "global-state-schema": {
      "num-byte-slice": number;
      "num-uint": number;
    },
    "local-state-schema": {
      "num-byte-slice": number;
      "num-uint": number;
    },
    "on-completion": string;
  };
  "inner-txns"?: TransactionResponse[];
  "payment-transaction"?: {
    amount: number;
    "close-amount": number;
    "close-remainder-to": string;
    receiver: string;
  };
  "asset-transfer-transaction": {
    "asset-id": number;
    amount: number;
    receiver: string;
    "close-amount": number;
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
  "asset-freeze-transaction": {
    address: string;
    "asset-id": number;
    "new-freeze-status": boolean;
  };
  fee: number;
  "round-time": number;
  "first-valid": number;
  "last-valid": number;
  timestamp: number;
  note: string;
  signature: {
    logicsig: {
      args?: [];
      logic: string | null;
      "multisig-signature"?: {};
    };
    multisig: {
      subsignature?: multisigSubsig[],
      threshold?: number;
      version?: number;
    };
    sig?: string;
  };
};

export interface IBlockRewards {
  "fee-sink": string;
  "rewards-calculation-round": number;
  "rewards-level": number;
  "rewards-pool": string;
  "rewards-rate": number;
  "rewards-residue": number;
}

export type IBlockResponse = {
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
  items: IBlockResponse[];
}

export interface ISupply {
  current_round: number;
  "online-money": string;
}
