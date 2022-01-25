import { commonIntercepts, interceptTxsOnTxsPage, stubWebSocketToInSync } from "../support/utils"

describe('Transactions Page', () => {
  beforeEach(() => {

    commonIntercepts();

    interceptTxsOnTxsPage();
    stubWebSocketToInSync();

    cy.visit('/transactions')

    cy.url().should('contain', 'http://localhost:3000/transactions')
  })

  it('displays transactions header text correctly', () => {
    cy.get('[class*="Breadcrumbs_pageTitle"]').should('have.text', 'Transactions')
    cy.get('[class*="Breadcrumbs_pageTitle"] + div').should('have.text', 'Home / Transactions')
  })

  it('clicking txs list tx id navigates to tx page', () => {
    cy.get('tbody tr:first-child td:first-child a').should('have.text', 'OVVVKO...RONF6A')
    cy.get('tbody tr:first-child td:first-child a').click()
    cy.url().should('include', '/tx/')
  })

  it('clicking txs list block number navigates to block page', () => {
    cy.get('tbody tr').eq(0).children().eq(1).children("a").as('blocknum')
    cy.get('@blocknum').should('have.text', '18,788,980')
    cy.get('@blocknum').click()
    cy.url().should('include', '/block/')
  })

  it('clicking txs list From navigates to address page', () => {
    cy.get('tbody tr').eq(0).children().eq(3).children("a").as('from')
    cy.get('@from').should('have.text', 'ZW3ISE...67W754')
    cy.get('@from').click({force:true})
    cy.wait(500).url().should('include', '/address/')
  })

  it('clicking txs list To navigates to address page', () => {
    cy.get('tbody tr').eq(0).children().eq(4).children("a").as('to')
    cy.get('@to').should('have.text', '6JUDNT...4NM6UE')
    cy.get('@to').click({force:true})
    cy.wait(500).url().should('include', '/address/')
  })

  it('shows To N/A for App Call Type transaction', () => {
    cy.get('tbody tr').eq(2).children().eq(2).should('have.text', 'App Call')
    cy.get('tbody tr').eq(2).children().eq(4).should('have.text', 'N/A')
  })

  it('shows page=1 on the url if not provided in url param', () => {
    cy.wait(500).url().should('contain', '?page=1')
  })

  it('shows page 3 on table pagination if page=3 is provided in url param', () => {
    cy.visit('/transactions', {
      qs: {
        page: 3
      }
    })
    cy.wait(500).url().should('contain', '?page=3')
    cy.get('[class*="CustomTable_page-input"]').should('have.value', 3)
  })

  it('shows page 1 on url param if user clicks table pagination button <<', () => {
    cy.visit('/transactions', {
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
    cy.visit('/transactions', {
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
    cy.url().should('contain', '?page=1')
    cy.get('[class*="CustomTable_page-input"]').should('have.value', 1)
    cy.wait(1000).get('[class*="CustomTable_page-input"]').type('{selectall}3').blur().invoke('val')
      .then(val=> {    
        cy.wait(500).url().should('contain', `?page=3`)
      })
  })

  it('shows page 2 on url param if user clicks table pagination button > on page 1', () => {
    cy.url().should('contain', '?page=1')
    cy.get('[class*="CustomTable_pagination"]').children().eq(3).click()
    cy.get('[class*="CustomTable_page-input"]').should('have.value', 2)
    cy.wait(500).url().should('contain', `?page=2`)
  })

  it('shows page 107100 on url param if user clicks table pagination button >>', () => {
    cy.url().should('contain', '?page=1')
    cy.get('[class*="CustomTable_pagination"]').children().eq(4).click()
    cy.get('[class*="CustomTable_page-input"]').should('have.value', 107100)
    cy.wait(500).url().should('contain', `?page=107100`)
  })

})
