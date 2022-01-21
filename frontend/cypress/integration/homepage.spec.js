describe('Home Page', () => {
  beforeEach(() => {
    // Cypress starts out with a blank slate for each test
    // so we must tell it to visit our website with the `cy.visit()` command.
    // Since we want to visit the same URL at the start of all our tests,
    // we include it in our beforeEach function so that it runs before each test
    cy.visit('http://localhost:3000/')

    cy.url().should('equal', 'http://localhost:3000/')
  })

  it('displays stats cards by default', () => {
    cy.get('*[class*="statscard"]').should('have.length', 5)

    cy.get('*[class*="statscard"]').first().should('have.text', 'Latest Round')
    cy.get('*[class*="statscard"]').last().should('have.text', 'Algo Price')
  })
})
