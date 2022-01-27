import { commonIntercepts, stubWebSocketToInSync } from "../support/utils"

describe('Asset Page', () => {
  beforeEach(() => {

    commonIntercepts();

    stubWebSocketToInSync();

    cy.visit('/asset/163650')

    cy.wait('@getAsset')

    cy.url().should('contain', 'http://localhost:3000/asset/163650')
  })

  it('displays asset info correctly', () => {
    cy.get('[class*="Block_block-table"] tbody tr').eq(0).children().eq(0).should('have.text', 'Asset ID')
    cy.get('[class*="Block_block-table"] tbody tr').eq(0).children().eq(1).should('have.text', '163650')
    cy.get('[class*="Block_block-table"] tbody tr').eq(1).children().eq(0).should('have.text', 'Asset Name')
    cy.get('[class*="Block_block-table"] tbody tr').eq(1).children().eq(1).should('have.text', 'Asia Reserve Currency Coin')
    cy.get('[class*="Block_block-table"] tbody tr').eq(2).children().eq(0).should('have.text', 'URL')
    cy.get('[class*="Block_block-table"] tbody tr').eq(2).children().eq(1).should('have.text', 'https://arcc.io')
    cy.get('[class*="Block_block-table"] tbody tr').eq(3).children().eq(0).should('have.text', 'Decimals')
    cy.get('[class*="Block_block-table"] tbody tr').eq(3).children().eq(1).should('have.text', '6')
    cy.get('[class*="Block_block-table"] tbody tr').eq(4).children().eq(0).should('have.text', 'Creator')
    cy.get('[class*="Block_block-table"] tbody tr').eq(4).children().eq(1).should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY')
    cy.get('[class*="Block_block-table"] tbody tr').eq(5).children().eq(0).should('have.text', 'Manager')
    cy.get('[class*="Block_block-table"] tbody tr').eq(5).children().eq(1).should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY')
    cy.get('[class*="Block_block-table"] tbody tr').eq(6).children().eq(0).should('have.text', 'Reserve Account')
    cy.get('[class*="Block_block-table"] tbody tr').eq(6).children().eq(1).should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY')
    cy.get('[class*="Block_block-table"] tbody tr').eq(7).children().eq(0).should('have.text', 'Freeze Account')
    cy.get('[class*="Block_block-table"] tbody tr').eq(7).children().eq(1).should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY')
    cy.get('[class*="Block_block-table"] tbody tr').eq(8).children().eq(0).should('have.text', 'Total Supply')
    cy.get('[class*="Block_block-table"] tbody tr').eq(8).children().eq(1).should('have.text', '88,616,203,378.51 ARCC')
  })

  it('check URL is set to asset\'s website for a new tab', () => {
    cy.get('[class*="Block_block-table"] tbody tr').eq(2).children().eq(1).children('a').as('url')
    cy.get('@url').should('have.text', 'https://arcc.io')
      .should('have.attr', 'href', 'https://arcc.io')
      .should('have.attr', 'target', '_blank')
      .should('have.attr', 'rel', 'noopener noreferrer')
  })

  it('clicking Creator navigates to address page', () => {
    cy.get('[class*="Block_block-table"] tbody tr').eq(4).children().eq(1).children('a').as('creator')
    cy.get('@creator').should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY')
    cy.get('@creator').click()
    cy.url().should('include', '/address/')
  })
  
  it('clicking Manager navigates to address page', () => {
    cy.get('[class*="Block_block-table"] tbody tr').eq(5).children().eq(1).children('a').as('manager')
    cy.get('@manager').should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY')
    cy.get('@manager').click()
    cy.url().should('include', '/address/')
  })

  it('clicking Reserve Account navigates to address page', () => {
    cy.get('[class*="Block_block-table"] tbody tr').eq(6).children().eq(1).children('a').as('reserve')
    cy.get('@reserve').should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY')
    cy.get('@reserve').click()
    cy.url().should('include', '/address/')
  })

  it('clicking Freeze Account navigates to address page', () => {
    cy.get('[class*="Block_block-table"] tbody tr').eq(7).children().eq(1).children('a').as('freeze')
    cy.get('@freeze').should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY')
    cy.get('@freeze').click()
    cy.url().should('include', '/address/')
  })
})
