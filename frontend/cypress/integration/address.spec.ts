import { backend_url, commonIntercepts, stubWebSocketToInSync } from '../support/utils';
const addressWappsFixture = '../fixtures/address/address_w_created_app.json';

describe('Address Page', () => {
  beforeEach(() => {
    commonIntercepts();

    stubWebSocketToInSync();

    cy.visit('/address/ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY');

    cy.wait('@getAddrTxs');

    cy.url().should('contain', 'http://localhost:3000/address/ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY');
  });

  it('shows page=1 on the url if not provided in url param', () => {
    cy.url().should('contain', '?page=1');
  });

  it('displays stats cards text correctly', () => {
    cy.get('*[class*="statscard"]').eq(0).children().first().should('have.text', 'Balance');
    cy.get('*[class*="statscard"]').eq(0).children().last().should('contain', '118.265216');
    cy.get('*[class*="statscard"]').eq(1).children().first().should('have.text', 'Rewards');
    cy.get('*[class*="statscard"]').eq(1).children().last().should('contain', '12.10355');
    cy.get('*[class*="statscard"]').eq(2).children().first().should('have.text', 'Pending Rewards');
    cy.get('*[class*="statscard"]').eq(2).children().last().should('contain', '2.706525');
    cy.get('*[class*="statscard"]').eq(3).children().first().children('h5').should('have.text', 'Status');
    cy.get('*[class*="statscard"]').eq(3).children().last().should('have.text', 'Offline');
  });

  /* Assets Table */
  it('clicking txs list tx id navigates to tx page', () => {
    cy.get('[class*="Block_table-wrapper"]').eq(0).find('tbody tr:first-child td:first-child a').as('asaid');
    cy.get('@asaid').should('have.text', '163650');
    cy.get('@asaid').click({ force: true });
    cy.wait(500).url().should('include', '/address/');
  });

  it('shows current address as Creator and is not clickable', () => {
    cy.get('[class*="Block_table-wrapper"]').eq(0).find('tbody tr:first-child td').eq(1).as('creator');
    cy.get('@creator').should('have.text', 'ARCC3T...GNMPSY');
    cy.get('@creator').find('a').should('have.lengthOf', 0);
  });

  it('shows asset info correctly', () => {
    cy.get('[class*="Block_table-header"]').eq(0).should('contain', 'Assets');
    cy.get('[class*="Block_table-wrapper"]').eq(0).find('tbody tr:first-child td').eq(2).should('have.text', '12.503738 ARCC');
    cy.get('[class*="Block_table-wrapper"]').eq(0).find('tbody tr:first-child td').eq(4).should('have.text', 'false');
  });

  /* Created Assets Table */
  it('shows created asset info correctly', () => {
    cy.get('[class*="Block_table-header"]').eq(1).should('contain', 'Created Assets');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(0).should('have.text', '163650');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(1).should('have.text', 'Asia Reserve Currency Coin');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(2).should('have.text', 'ARCC3T...GNMPSY');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(3).should('have.text', 'ARCC3T...GNMPSY');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(4).should('have.text', 'ARCC3T...GNMPSY');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(5).should('have.text', 'ARCC3T...GNMPSY');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(6).should('have.text', '6');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(7).should('have.text', '88,616,203,378.51 ARCC');
  });

  /* Created Apps Table */
  it('shows created apps info correctly', () => {
    cy.visit('/address/MIZ3P3RZEXYT4VHRRT6K5EBYQKP24SLHBPXMADKKDZ3VCLVXOOKUACN42E');
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/accounts/*?page=*&limit=10&order=desc`,
      },
      {
        fixture: addressWappsFixture,
      },
    ).as('getAddr2');
    cy.wait('@getAddr2');

    cy.get('[class*="Block_table-header"]').eq(1).should('contain', 'Created Apps');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(0).should('have.text', '62368684');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(1).should('have.text', '19,193,933');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(2).should('contain', '# byte-slice: 0');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(2).should('contain', '# uint: 0');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(3).should('contain', '# byte-slice: 0');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(3).should('contain', '# uint: 16');
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('tbody tr:first-child td').eq(4).should('have.text', 'false');
  });

  /* Transaction Table */
  it('clicking txs list tx id navigates to tx page', () => {
    cy.get('[class*="Block_table-header"]').eq(2).should('contain', 'Transactions');
    cy.get('[class*="Block_table-wrapper"]').eq(2).find('tbody tr:first-child td:first-child a').as('txid');
    cy.get('@txid').should('have.text', 'TTB72Z...XU7KLQ');
    cy.get('@txid').click({ force: true });
    cy.wait(500).url().should('include', '/tx/');
  });

  it('clicking txs list block number navigates to block page', () => {
    cy.get('[class*="Block_table-wrapper"]').eq(2).find('tbody tr').eq(0).children().eq(1).children('a').as('blocknum');
    cy.get('@blocknum').should('have.text', '6,640,016');
    cy.get('@blocknum').click({ force: true });
    cy.url().should('include', '/block/');
  });

  it('clicking txs list From navigates to address page', () => {
    cy.get('[class*="Block_table-wrapper"]').eq(2).find('tbody tr').eq(2).children().eq(3).children().as('from');
    cy.get('@from').should('have.text', 'DEKUZQ...LGEULQ');
    cy.get('@from').click({ force: true });
    cy.url().should('include', '/address/');
  });

  it('shows current address as To and is not clickable', () => {
    cy.get('[class*="Block_table-wrapper"]').eq(2).find('tbody tr').eq(2).children().eq(4).children().as('to');
    cy.get('@to').should('have.text', 'ARCC3T...GNMPSY');
    cy.get('@to').find('a').should('have.lengthOf', 0);
  });

  /* URL Page Param Effect on Tables */
  it('shows page 3 on url param if user inputs 3 as page number', () => {
    cy.wait(500).url().should('contain', '?page=1');
    cy.get('[class*="Block_table-wrapper"]').eq(0).find('[class*="CustomTable_page-input"]').should('have.value', 1);
    cy.get('[class*="Block_table-wrapper"]').eq(1).find('[class*="CustomTable_page-input"]').should('have.value', 1);
    cy.get('[class*="Block_table-wrapper"]').eq(2).find('[class*="CustomTable_page-input"]').should('have.value', 1);
    cy.wait(1000)
      .get('[class*="Block_table-wrapper"]')
      .eq(2)
      .find('[class*="CustomTable_page-input"]')
      .type('{selectall}3')
      .blur()
      .invoke('val')
      .then(val => {
        cy.wait(500).url().should('contain', `?page=3`);
        cy.get('[class*="Block_table-wrapper"]').eq(0).find('[class*="CustomTable_page-input"]').should('have.value', 1);
        cy.get('[class*="Block_table-wrapper"]').eq(1).find('[class*="CustomTable_page-input"]').should('have.value', 1);
      });
  });
});
