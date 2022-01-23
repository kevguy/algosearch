const supplyFixture = "../fixtures/supply.json";
const currentTxFixture = "../fixtures/tx/tx_pay_single.json";
const blocksFixture = "../fixtures/blocks/blocks_pay_txs.json";
const txsFixture = "../fixtures/txs/txs_pay.json";

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
      url: `${backend_url}/v1/rounds?latest_blk=*&limit=15&page=1&order=desc`,
    },
    {
      fixture: blocksFixture
    }
  ).as('getBlocks')
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
      url: `${backend_url}/v1/transactions?latest_txn=*&page=1&limit=15&order=desc`,
    },
    {
      fixture: txsFixture
    }
  ).as('getTxs')
}
