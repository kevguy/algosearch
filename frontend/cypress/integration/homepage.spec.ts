import { Server, WebSocket } from "mock-socket";

const supplyFixture = "../fixtures/supply.json";
const currentTxFixture = "../fixtures/tx/tx_pay_single.json";
const blocksFixture = "../fixtures/blocks/blocks_pay_txs.json";
const txsFixture = "../fixtures/txs/txs_pay.json";

const backend_url = "http://localhost:5000";

window.WebSocket = WebSocket; // Here we stub out the window object

let mockSocket;
let mockServer;
describe('Home Page', () => {
  beforeEach(() => {
    // Cypress starts out with a blank slate for each test
    // so we must tell it to visit our website with the `cy.visit()` command.
    // Since we want to visit the same URL at the start of all our tests,
    // we include it in our beforeEach function so that it runs before each test

    cy.intercept(
      {
        method: 'GET',
        url: `https://api.coingecko.com/api/v3/simple/price?ids=algorand&vs_currencies=usd`,
      },
      {
        statusCode: 200,
        body: {"algorand":{"usd":0.908584}}
      }
    )

    cy.intercept(
      {
        method: 'GET',
        url: `https://metricsapi.algorand.foundation/v1/supply/circulating?unit=algo`,
      },
      {
        statusCode: 200,
        body: 6530074763.318415
      }
    )

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

    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/rounds?latest_blk=*&page=1&limit=10&order=desc`,
      },
      {
        fixture: blocksFixture
      }
    ).as('getLatestBlocks')

    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/transactions?latest_txn=*&page=1&limit=10&order=desc`,
      },
      {
        fixture: txsFixture
      }
    ).as('getLatestTxs')

    cy.visit('/')
    cy.url().should('equal', 'http://localhost:3000/')
  })

  it('displays home header text correctly', () => {
    cy.get('*[class*="HomeHeader_content"] h1 span').should('have.text', 'Algorand Block Explorer')
    cy.get('*[class*="HomeHeader_desc"] > span').should('have.text', 'Open-source block explorer for Algorand')
  })

  it('shows in sync when it is in sync', () => {
    // TODO -> stub sync status
    cy.get('*[class*="sync-status"]').should('have.text', 'in sync')
  })

  it('shows out of sync when it is out of sync', () => {
    // TODO -> stub sync status
    // also check for out of sync by {num} blocks
    cy.get('*[class*="sync-status"]').should('have.text', 'out of sync')
  })

  it('displays stats cards by default', () => {
    cy.get('*[class*="statscard"]').should('have.length', 5)

    cy.get('*[class*="statscard"]').eq(0).should('have.text', 'Latest Round')
    cy.get('*[class*="statscard"]').eq(1).should('have.text', 'Online Stake')
    cy.get('*[class*="statscard"]').eq(2).should('have.text', 'Circulating Supply')
    cy.get('*[class*="statscard"]').eq(3).should('have.text', 'Block Time')
    cy.get('*[class*="statscard"]').eq(4).should('have.text', 'Algo Price')
  })

  it('clicking blocks list block number navigates to block page', () => {
    cy.wait('@getLatestBlocks', {timeout: 15000})
    cy.get('*[class*="BlockTable_block-row"]:first-child [class*="block-id"]').click()

    cy.url().should('include', '/block/')
  })

  it('clicking blocks list proposer navigates to address page', () => {
    cy.wait('@getLatestBlocks', {timeout: 15000})
    cy.get('*[class*="BlockTable_block-row"]:first-child [class*="proposer"]').click()

    cy.url().should('include', '/address/')
  })

  it('clicking transactions list tx id navigates to tx page', () => {
    cy.wait('@getLatestTxs', {timeout: 15000})
    cy.get('*[class*="TransactionTable_transaction-row"]:first-child [class*="transaction-id"]').click()

    cy.url().should('include', '/tx/')
  })

  it('clicking transactions list From navigates to address page', () => {
    cy.wait('@getLatestTxs', {timeout: 15000})
    cy.get('*[class*="TransactionTable_transaction-row"]:first-child [class*="TransactionTable_relevant-accounts"] span:first-child').click()

    cy.url().should('include', '/address/')
  })

  it('clicking transactions list To navigates to address page', () => {
    cy.wait('@getLatestTxs', {timeout: 15000})
    cy.get('*[class*="TransactionTable_transaction-row"]:first-child [class*="TransactionTable_relevant-accounts"]').children().eq(1).click()

    cy.url().should('include', '/address/')
  })
})
