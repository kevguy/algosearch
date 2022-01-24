import { commonIntercepts, interceptBlocksOnBlocksPage, stubWebSocketToInSync } from "../support/utils"

describe('Blocks Page', () => {
  beforeEach(() => {

    commonIntercepts();

    interceptBlocksOnBlocksPage();
    stubWebSocketToInSync();

    cy.visit('/blocks')

    cy.url().should('contain', 'http://localhost:3000/blocks')
  })

  it('displays blocks header text correctly', () => {
    cy.get('[class*="Breadcrumbs_pageTitle"]').should('have.text', 'Blocks')
    cy.get('[class*="Breadcrumbs_pageTitle"] + div').should('have.text', 'Home / Blocks')
  })

  it('clicking blocks list block number navigates to block page', () => {
    cy.get('tbody tr:first-child a').should('have.text', '18,788,980')
    cy.get('tbody tr:first-child a').click()
    cy.url().should('include', '/block/')
  })

  it('clicking blocks list proposer navigates to address page', () => {
    cy.get('tbody tr').eq(3).children().eq(1).children("a").should('have.text', 'REMF54...NCHH6M')
    cy.get('tbody tr').eq(3).children().eq(1).children("a").click()
    cy.url().should('include', '/address/')
  })

  it('displays stats cards text correctly', () => {
    cy.get('*[class*="statscard"]').eq(0).children().first().should('have.text', 'Latest Block')
    cy.get('*[class*="statscard"]').eq(0).children().last().should('have.text', '18,788,980')
    cy.get('*[class*="statscard"]').eq(1).children().first().should('have.text', 'Block Time')
    cy.get('*[class*="statscard"]').eq(1).children().last().should('have.text', '4.375 seconds')
  })

  it('shows page=1 on the url if not provided in url param', () => {
    cy.wait(500).url().should('contain', '?page=1')
  })

  it('shows page 3 on table pagination if page=3 is provided in url param', () => {
    cy.visit('/blocks', {
      qs: {
        page: 3
      }
    })
    cy.wait(500).url().should('contain', '?page=3')
    cy.get('[class*="CustomTable_page-input"]').should('have.value', 3)
  })

  it('shows page 1 on url param if user clicks table pagination button <<', () => {
    cy.visit('/blocks', {
      qs: {
        page: 10
      }
    })
    cy.wait(500).url().should('contain', '?page=10')
    cy.get('[class*="CustomTable_pagination"]').children().eq(0).click()
    cy.get('[class*="CustomTable_page-input"]').should('have.value', 1)
    cy.wait(500).url().should('contain', `?page=1`)
  })

  it('shows page 9 on url param if user clicks table pagination button < on page 10', () => {
    cy.visit('/blocks', {
      qs: {
        page: 10
      }
    })
    cy.wait(500).url().should('contain', '?page=10')
    cy.get('[class*="CustomTable_pagination"]').children().eq(1).click()
    cy.get('[class*="CustomTable_page-input"]').should('have.value', 9)
    cy.wait(500).url().should('contain', `?page=9`)
  })

  it('shows page 3 on url param if user inputs 3 as page number', () => {
    cy.wait(500).url().should('contain', '?page=1')
    cy.get('[class*="CustomTable_page-input"]').should('have.value', 1)
    cy.wait(1000).get('[class*="CustomTable_page-input"]').type('{selectall}3').blur().invoke('val')
      .then(val=> {    
        cy.wait(500).url().should('contain', `?page=3`)
      })
  })

  it('shows page 2 on url param if user clicks table pagination button > on page 1', () => {
    cy.wait(500).url().should('contain', '?page=1')
    cy.get('[class*="CustomTable_pagination"]').children().eq(3).click()
    cy.get('[class*="CustomTable_page-input"]').should('have.value', 2)
    cy.wait(500).url().should('contain', `?page=2`)
  })

  it('shows page 453371 on url param if user clicks table pagination button >>', () => {
    cy.wait(500).url().should('contain', '?page=1')
    cy.get('[class*="CustomTable_pagination"]').children().eq(4).click()
    cy.get('[class*="CustomTable_page-input"]').should('have.value', 453371)
    cy.wait(500).url().should('contain', `?page=453371`)
  })

})
