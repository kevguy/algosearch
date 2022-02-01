import { backend_url, commonIntercepts, interceptBlocks, interceptTxs, stubWebSocketToInSync, stubWebSocketToOutOfSync } from '../support/utils';
const searchTxFixture = '../fixtures/search/search_tx.json';
const searchAddrFixture = '../fixtures/search/search_address.json';
const searchBlockFixture = '../fixtures/search/search_block.json';

describe('Home Page', () => {
  beforeEach(() => {
    cy.intercept(
      {
        method: 'GET',
        url: `https://api.coingecko.com/api/v3/simple/price?ids=algorand&vs_currencies=usd`,
      },
      {
        statusCode: 200,
        body: { algorand: { usd: 0.908584 } },
      },
    );

    cy.intercept(
      {
        method: 'GET',
        url: `https://metricsapi.algorand.foundation/v1/supply/circulating?unit=algo`,
      },
      {
        statusCode: 200,
        body: 6530074763.318415,
      },
    );

    commonIntercepts();

    interceptBlocks();

    interceptTxs();
  });

  it('displays home header text correctly', () => {
    cy.visit('/');
    cy.url().should('equal', 'http://localhost:3000/');
    cy.get('*[class*="HomeHeader_content"] h1 span').should('have.text', 'Algorand Block Explorer');
    cy.get('*[class*="HomeHeader_desc"] > span').should('have.text', 'Open-source block explorer for Algorand');
  });

  it('shows in sync when it is in sync', () => {
    stubWebSocketToInSync();
    cy.visit('/');

    cy.get('*[class*="sync-status"]').should('have.text', 'in sync');
  });

  it('shows out of sync when it is out of sync', () => {
    stubWebSocketToOutOfSync();
    cy.visit('/');

    cy.get('*[class*="sync-status"]').should('have.text', 'out of sync by 14,529,128 blocks');
  });

  it('displays stats cards text correctly', () => {
    stubWebSocketToInSync();
    cy.visit('/');
    cy.wait('@getLatestTxs');
    cy.wait('@getLatestBlocks');

    cy.get('*[class*="statscard"]').should('have.length', 5);

    cy.get('*[class*="statscard"]').eq(0).children().first().should('have.text', 'Latest Round');
    cy.get('*[class*="statscard"]').eq(0).children().last().should('have.text', '18,788,980');
    cy.get('*[class*="statscard"]').eq(1).children().first().should('have.text', 'Online Stake');
    cy.get('*[class*="statscard"]').eq(1).children().last().should('include.text', '2,099,856,660.7');
    cy.get('*[class*="statscard"]').eq(2).children().first().should('have.text', 'Circulating Supply');
    cy.get('*[class*="statscard"]').eq(2).children().last().should('include.text', '6,530,074,763.318415');
    cy.get('*[class*="statscard"]').eq(3).children().first().should('have.text', 'Block Time');
    cy.get('*[class*="statscard"]').eq(3).children().last().should('have.text', '4.375 seconds');
    cy.get('*[class*="statscard"]').eq(4).children().first().should('have.text', 'Algo Price');
    cy.get('*[class*="statscard"]').eq(4).children().last().should('have.text', '$0.908584');
  });

  it('clicking blocks list block number navigates to block page', () => {
    stubWebSocketToInSync();
    cy.visit('/');
    cy.wait('@getLatestTxs');
    cy.wait('@getLatestBlocks');

    cy.get('*[class*="BlockTable_block-row"]:first-child [class*="block-id"]').click();

    cy.url().should('include', '/block/');
  });

  it('clicking blocks list proposer navigates to address page', () => {
    stubWebSocketToInSync();
    cy.visit('/');
    cy.wait('@getLatestTxs');
    cy.wait('@getLatestBlocks');

    cy.get('*[class*="BlockTable_block-row"]:first-child [class*="proposer"]').click();

    cy.url().should('include', '/address/');
  });

  it('clicking transactions list tx id navigates to tx page', () => {
    stubWebSocketToInSync();
    cy.visit('/');
    cy.wait('@getLatestTxs');
    cy.wait('@getLatestBlocks');

    cy.get('*[class*="TransactionTable_transaction-row"]:first-child [class*="transaction-id"]').click();

    cy.wait(5000).url().should('include', '/tx/');
  });

  it('clicking transactions list From navigates to address page', () => {
    stubWebSocketToInSync();
    cy.visit('/');
    cy.wait('@getLatestTxs');
    cy.wait('@getLatestBlocks');

    cy.get('*[class*="TransactionTable_transaction-row"]:first-child [class*="TransactionTable_relevant-accounts"] span:first-child').click();

    cy.url().should('include', '/address/');
  });

  it('clicking transactions list To navigates to address page', () => {
    stubWebSocketToInSync();
    cy.visit('/');
    cy.wait('@getLatestTxs');
    cy.wait('@getLatestBlocks');

    cy.get('*[class*="TransactionTable_transaction-row"]:first-child [class*="TransactionTable_relevant-accounts"]').children().eq(1).click();

    cy.wait(500).url().should('include', '/address/');
  });

  /* Test Header Search Bar */
  it('searching an address on header search bar navigates to address page', () => {
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/search?key=ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY`,
      },
      {
        fixture: searchAddrFixture,
      },
    ).as('searchAddr');

    stubWebSocketToInSync();
    cy.visit('/');

    cy.get('input[type="search"]').type('ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY');
    cy.get('input[type="search"]').siblings('button').click();
    cy.wait('@searchAddr');

    cy.wait(500).url().should('include', '/address/ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY');
  });

  it('searching a tx id on header search bar navigates to tx page', () => {
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/search?key=NTIU26TLJ6XMMBV6YQJB6SUPG5FBKCMHG2EQ5R5AGJDQ7OXK7PKQ`,
      },
      {
        fixture: searchTxFixture,
      },
    ).as('searchTx');

    stubWebSocketToInSync();
    cy.visit('/');

    cy.get('input[type="search"]').type('NTIU26TLJ6XMMBV6YQJB6SUPG5FBKCMHG2EQ5R5AGJDQ7OXK7PKQ');
    cy.get('input[type="search"]').siblings('button').click();
    cy.wait('@searchTx');

    cy.wait(500).url().should('include', '/tx/NTIU26TLJ6XMMBV6YQJB6SUPG5FBKCMHG2EQ5R5AGJDQ7OXK7PKQ');
  });

  it('searching a block id on header search bar navigates to block page', () => {
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/search?key=4259852`,
      },
      {
        fixture: searchBlockFixture,
      },
    ).as('searchBlock');

    stubWebSocketToInSync();
    cy.visit('/');

    cy.get('input[type="search"]').type('4259852');
    cy.get('input[type="search"]').siblings('button').click();
    cy.wait('@searchBlock');

    cy.wait(500).url().should('include', '/block/4259852');
  });
});
