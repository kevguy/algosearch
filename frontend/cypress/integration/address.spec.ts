import { commonIntercepts, stubWebSocketToInSync } from "../support/utils"

describe('Address Page', () => {
  beforeEach(() => {

    commonIntercepts();

    stubWebSocketToInSync();

    cy.visit('/address/ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY')

    cy.wait('@getAddrTxs')

    cy.url().should('contain', 'http://localhost:3000/address/ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY')
  })

  it('shows page=1 on the url if not provided in url param', () => {
    cy.url().should('contain', '?page=1')
  })

  it('displays stats cards text correctly', () => {
    cy.get('*[class*="statscard"]').eq(0).children().first().should('have.text', 'Balance')
    cy.get('*[class*="statscard"]').eq(0).children().last().should('contain', '115.558691')
    cy.get('*[class*="statscard"]').eq(1).children().first().should('have.text', 'Rewards')
    cy.get('*[class*="statscard"]').eq(1).children().last().should('contain', '12.10355')
    cy.get('*[class*="statscard"]').eq(2).children().first().should('have.text', 'Pending Rewards')
    cy.get('*[class*="statscard"]').eq(2).children().last().should('contain', '2.706525')
    cy.get('*[class*="statscard"]').eq(3).children().first().should('have.text', 'Transactions')
    cy.get('*[class*="statscard"]').eq(3).children().last().should('contain', '58')
    cy.get('*[class*="statscard"]').eq(4).children().first().children("h5").should('have.text', 'Status')
    cy.get('*[class*="statscard"]').eq(4).children().last().should('have.text', 'Offline')
  })

  it('clicking txs list tx id navigates to tx page', () => {
    cy.get('tbody tr:first-child td:first-child a').as('txid')
    cy.get('@txid').should('have.text', 'TTB72Z...XU7KLQ')
    cy.get('@txid').click({force: true})
    cy.url().should('include', '/tx/')
  })

  it('clicking txs list block number navigates to block page', () => {
    cy.get('tbody tr').eq(0).children().eq(1).children("a").as('blocknum')
    cy.get('@blocknum').should('have.text', '6,640,016')
    cy.get('@blocknum').click({force: true})
    cy.url().should('include', '/block/')
  })

  it('clicking txs list From navigates to address page', () => {
    cy.get('tbody tr').eq(2).children().eq(3).children().as('from')
    cy.get('@from').should('have.text', 'DEKUZQ...LGEULQ')
    cy.get('@from').click({force: true})
    cy.url().should('include', '/address/')
  })

  it('shows current address as To and is not clickable', () => {
    cy.get('tbody tr').eq(2).children().eq(4).children().as('to')
    cy.get('@to').should('have.text', 'ARCC3T...GNMPSY')
  })
})
