import { commonIntercepts, stubWebSocketToInSync } from "../support/utils"

describe('Block Page', () => {
  beforeEach(() => {

    commonIntercepts();

    stubWebSocketToInSync();

    cy.visit('/block/18788980')

    cy.wait('@getBlock')

    cy.url().should('contain', 'http://localhost:3000/block/18788980')
  })

  it('displays block info correctly', () => {
    cy.get('[class*="Block_block-table"] tbody tr').eq(0).children().eq(0).should('have.text', 'Block')
    cy.get('[class*="Block_block-table"] tbody tr').eq(0).children().eq(1).should('have.text', '18788980')
    cy.get('[class*="Block_block-table"] tbody tr').eq(2).children().eq(0).should('have.text', 'Reward Rate')
    cy.get('[class*="Block_block-table"] tbody tr').eq(2).children().eq(1).should('contain', '52.000012')
    cy.get('[class*="Block_block-table"] tbody tr').eq(3).children().eq(0).should('have.text', 'Proposer')
    cy.get('[class*="Block_block-table"] tbody tr').eq(3).children().eq(1).should('have.text', 'REMF542E5ZFKS7SGSNHTYB255AUITEKHLAATWVPK3CY7TAFPT6GNNCHH6M')
    cy.get('[class*="Block_block-table"] tbody tr').eq(4).children().eq(0).should('have.text', 'Block Hash')
    cy.get('[class*="Block_block-table"] tbody tr').eq(4).children().eq(1).should('have.text', 'B2WNSa7XTUmY8fnq1sOnCU2kzouykni/fRU2w0WdIOE=')
    cy.get('[class*="Block_block-table"] tbody tr').eq(5).children().eq(0).should('have.text', 'Previous Block Hash')
    cy.get('[class*="Block_block-table"] tbody tr').eq(5).children().eq(1).should('have.text', 'xD23ffs8gvOvURZbgQtVtxKhQMCc5wtoohB0sq4FuBA=')
    cy.get('[class*="Block_block-table"] tbody tr').eq(6).children().eq(0).should('have.text', 'Seed')
    cy.get('[class*="Block_block-table"] tbody tr').eq(6).children().eq(1).should('have.text', '+oqnjMAeCd/RJBGZ9F4pQcNRWpM9VstN2z+HGAlvrgQ=')
    cy.get('[class*="Block_block-table"] tbody tr').eq(7).children().eq(0).should('have.text', 'Transactions')
    cy.get('[class*="Block_block-table"] tbody tr').eq(7).children().eq(1).should('have.text', '4')
  })

  it('shows page=1 on the url if not provided in url param', () => {
    cy.url().should('contain', '?page=1')
  })

  it('clicking txs list tx id navigates to tx page', () => {
    cy.get('[class*="Block_transactions-table"] tbody tr:first-child td:first-child a').as('txid')
    cy.get('@txid').should('have.text', '5FKLBB...3LM5KQ')
    cy.get('@txid').click({force: true})
    cy.url().should('include', '/tx/')
  })

  it('clicking txs list From navigates to address page', () => {
    cy.get('[class*="Block_transactions-table"] tbody tr:first-child').children().eq(2).children().as('from')
    cy.get('@from').should('have.text', 'ARCC3T...GNMPSY')
    cy.get('@from').click({force: true})
    cy.url().should('include', '/address/')
  })

  it('clicking txs list To navigates to address page', () => {
    cy.get('[class*="Block_transactions-table"] tbody tr:first-child').children().eq(3).children().as('to')
    cy.get('@to').should('have.text', 'DBOVYJ...FVB2G4')
    cy.get('@to').click({force: true})
    cy.url().should('include', '/address/')
  })
})
