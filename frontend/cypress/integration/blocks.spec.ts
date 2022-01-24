import { commonIntercepts, interceptBlocksOnBlocksPage } from "../support/utils"

describe('Blocks Page', () => {
  beforeEach(() => {

    commonIntercepts();

    interceptBlocksOnBlocksPage();

    cy.visit('/blocks')

    cy.url().should('contain', 'http://localhost:3000/blocks')
  })

  it('displays blocks header text correctly', () => {
    cy.get('[class*="Breadcrumbs_pageTitle"]').should('have.text', 'Blocks')
    cy.get('[class*="Breadcrumbs_pageTitle"] + div').should('have.text', 'Home / Blocks')
  })

  it('displays stats cards text correctly', () => {
    cy.get('*[class*="statscard"]').eq(0).should('have.text', 'Latest Block')
    cy.get('*[class*="statscard"]').eq(1).should('have.text', 'Block Time')
  })

})
