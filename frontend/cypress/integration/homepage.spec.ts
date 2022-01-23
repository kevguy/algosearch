import * as blockInSyncFixture from "../fixtures/block/block_18788980.json";
import * as blockOutOfSyncFixture from "../fixtures/block/block_4259852.json";
import { commonIntercepts, interceptBlocks, interceptTxs } from "./utils";

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

    commonIntercepts();

    interceptBlocks();

    interceptTxs();
  })

  it('displays home header text correctly', () => {
    cy.visit('/')
    cy.url().should('equal', 'http://localhost:3000/')
    cy.get('*[class*="HomeHeader_content"] h1 span').should('have.text', 'Algorand Block Explorer')
    cy.get('*[class*="HomeHeader_desc"] > span').should('have.text', 'Open-source block explorer for Algorand')
  })

  it('shows in sync when it is in sync', () => {
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
    cy.visit('/')

    cy.get('*[class*="sync-status"]').should('have.text', 'in sync')
  })

  it('shows out of sync when it is out of sync', () => {
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
    cy.visit('/')

    cy.get('*[class*="sync-status"]').should('have.text', 'out of sync by 14,529,128 blocks')
  })

  it('displays stats cards text correctly', () => {
    cy.visit('/')
    cy.url().should('equal', 'http://localhost:3000/')

    cy.get('*[class*="statscard"]').should('have.length', 5)

    cy.get('*[class*="statscard"]').eq(0).should('have.text', 'Latest Round')
    cy.get('*[class*="statscard"]').eq(1).should('have.text', 'Online Stake')
    cy.get('*[class*="statscard"]').eq(2).should('have.text', 'Circulating Supply')
    cy.get('*[class*="statscard"]').eq(3).should('have.text', 'Block Time')
    cy.get('*[class*="statscard"]').eq(4).should('have.text', 'Algo Price')
  })

  it('clicking blocks list block number navigates to block page', () => {
    cy.visit('/')
    cy.url().should('equal', 'http://localhost:3000/')
    cy.wait('@getLatestBlocks', {timeout: 15000})
    cy.get('*[class*="BlockTable_block-row"]:first-child [class*="block-id"]').click()

    cy.url().should('include', '/block/')
  })

  it('clicking blocks list proposer navigates to address page', () => {
    cy.visit('/')
    cy.url().should('equal', 'http://localhost:3000/')
    cy.wait('@getLatestBlocks', {timeout: 15000})
    cy.get('*[class*="BlockTable_block-row"]:first-child [class*="proposer"]').click()

    cy.url().should('include', '/address/')
  })

  it('clicking transactions list tx id navigates to tx page', () => {
    cy.visit('/')
    cy.url().should('equal', 'http://localhost:3000/')
    cy.wait('@getLatestTxs', {timeout: 15000})
    cy.get('*[class*="TransactionTable_transaction-row"]:first-child [class*="transaction-id"]').click()

    cy.url().should('include', '/tx/')
  })

  it('clicking transactions list From navigates to address page', () => {
    cy.visit('/')
    cy.url().should('equal', 'http://localhost:3000/')
    cy.wait('@getLatestTxs', {timeout: 15000})
    cy.get('*[class*="TransactionTable_transaction-row"]:first-child [class*="TransactionTable_relevant-accounts"] span:first-child').click()

    cy.url().should('include', '/address/')
  })

  it('clicking transactions list To navigates to address page', () => {
    cy.visit('/')
    cy.url().should('equal', 'http://localhost:3000/')
    cy.wait('@getLatestTxs', {timeout: 15000})
    cy.get('*[class*="TransactionTable_transaction-row"]:first-child [class*="TransactionTable_relevant-accounts"]').children().eq(1).click()

    cy.url().should('include', '/address/')
  })
})
