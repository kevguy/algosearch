const supplyFixture = "../fixtures/supply.json";
const currentTxFixture = "../fixtures/tx/tx_pay_single.json";
const addressFixture = "../fixtures/address/address_1.json";
const addressTxsFixture = "../fixtures/address/address_1_txs.json";
const assetFixture = "../fixtures/asset/asset_1.json";
const blocksFixture = "../fixtures/blocks/blocks_pay_txs.json";
const blocksMixedFixture = "../fixtures/blocks/blocks_mixed_txs.json";
import * as blockInSyncFixture from "../fixtures/block/block_18788980.json";
import * as blockOutOfSyncFixture from "../fixtures/block/block_4259852.json";
const txsMixedFixture = "../fixtures/txs/txs_mixed.json";
const txsFixture = "../fixtures/txs/txs_pay.json";
const txPaySingleFixture = "../fixtures/tx/tx_pay_single.json";

export const backend_url = "http://localhost:5000";

export const commonIntercepts = () => {
  cy.intercept(
    {
      method: 'GET',
      url: `${backend_url}/v1/current-txn`,
    },
    {
      fixture: currentTxFixture
    }
  )

  cy.intercept(
    {
      method: 'GET',
      url: `${backend_url}/v1/algod/ledger/supply`,
    },
    {
      fixture: supplyFixture
    }
  )

  // interceptNext()

  interceptAddr()
  interceptAddrTxs()
  interceptAsset()
  interceptTx()
}

export const interceptNext = () => {
  cy.intercept(
    {
      method: 'GET',
      url: 'http://localhost:3000/_next/static/development/_devPagesManifest.json',
    },
    {
      response: 200,
      headers: {
        "host": "localhost:3000",
        "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:96.0) Gecko/20100101 Firefox/96.0",
        "accept": "*/*",
        "accept-language": "en-US,en;q=0.5",
        "accept-encoding": "gzip, deflate",
        "connection": "keep-alive"
      },
      body: {
        "pages": [
          "/",
          "/404",
          "/_app",
          "/_document",
          "/address/[_address]",
          "/api/__coverage__",
          "/asset/[_asset]",
          "/block/[_block]",
          "/blocks",
          "/transactions",
          "/tx/TransactionDetails",
          "/tx/[_txid]"
        ]
      }
    }
  )
}

export const interceptAddr = () => {
  cy.intercept(
    {
      method: 'GET',
      url: `${backend_url}/v1/accounts/*?page=*&limit=10&order=desc`,
    },
    {
      fixture: addressFixture
    }
  ).as('getAddr')
}

export const interceptAddrTxs = () => {
  cy.intercept(
    {
      method: 'GET',
      url: `${backend_url}/v1/transactions/acct/*?page=*&limit=15`,
    },
    {
      fixture: addressTxsFixture
    }
  ).as('getAddrTxs')
}

export const interceptAsset = () => {
  cy.intercept(
    {
      method: 'GET',
      url: `${backend_url}/v1/algod/assets/*`,
    },
    {
      fixture: assetFixture
    }
  ).as('getAsset')
}

export const interceptBlocks = () => {
  cy.intercept(
    {
      method: 'GET',
      url: `${backend_url}/v1/rounds?latest_blk=*&page=1&limit=10&order=desc`,
    },
    {
      fixture: blocksFixture
    }
  ).as('getLatestBlocks')
}

export const interceptBlocksOnBlocksPage = () => {
  cy.intercept(
    {
      method: 'GET',
      url: `${backend_url}/v1/rounds?latest_blk=*&limit=15&page=*&order=desc`,
    },
    {
      fixture: blocksMixedFixture
    }
  ).as('getBlocks')
}

export const interceptTx = () => {
  cy.intercept(
    {
      method: 'GET',
      url: `${backend_url}/v1/transactions/*`
    },
    {
      fixture: txPaySingleFixture
    }
  ).as('getTx')
}

export const interceptTxs = () => {
  cy.intercept(
    {
      method: 'GET',
      url: `${backend_url}/v1/transactions?latest_txn=*&page=1&limit=10&order=desc`,
    },
    {
      fixture: txsFixture
    }
  ).as('getLatestTxs')
}

export const interceptTxsOnTxsPage = () => {
  cy.intercept(
    {
      method: 'GET',
      url: `${backend_url}/v1/transactions?latest_txn=*&page=*&limit=15&order=desc`,
    },
    {
      fixture: txsMixedFixture
    }
  ).as('getTxs')
}

export const stubWebSocketToInSync = () => {
  Cypress.on("window:before:load", win => {
    (win as any).useWebSocketLibHook = () => ({
      sendMessage: () => {},
      sendJsonMessage: () => {},
      lastMessage: {},
      lastJsonMessage: {
        account_ids: [],
        app_ids: null,
        asset_ids: null,
        avg_block_txn_speed: 4.375,
        block: blockInSyncFixture,
        transaction_ids: [],
      },
      readyState: 1,
      getWebSocket: () => {},
    });
  })
}

export const stubWebSocketToOutOfSync = () => {
  Cypress.on("window:before:load", win => {
    (win as any).useWebSocketLibHook = () => ({
      sendMessage: () => {},
      sendJsonMessage: () => {},
      lastMessage: {},
      lastJsonMessage: {
        account_ids: [],
        app_ids: null,
        asset_ids: null,
        avg_block_txn_speed: 4.375,
        block: blockOutOfSyncFixture,
        transaction_ids: [],
      },
      readyState: 1,
      getWebSocket: () => {},
    });
  })
}
