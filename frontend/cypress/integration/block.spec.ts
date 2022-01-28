import { backend_url, commonIntercepts, stubWebSocketToInSync } from "../support/utils"
const blockFixture = "../fixtures/block/block_18788980.json";

describe('Block Page', () => {
  beforeEach(() => {

    commonIntercepts();
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/algod/rounds/*`,
      },
      {
        fixture: blockFixture
      }
    ).as('getInSyncBlock')

    stubWebSocketToInSync();

    cy.visit('/block/18788980')

    cy.wait('@getInSyncBlock')

    cy.url().should('contain', 'http://localhost:3000/block/18788980')
  })

  it('displays block info correctly', () => {
    cy.get('[class*="Block_block-table"] tbody tr').eq(0).children().eq(0).should('have.text', 'Block')
    cy.get('[class*="Block_block-table"] tbody tr').eq(0).children().eq(1).should('have.text', '18788980')
    cy.get('[class*="Block_block-table"] tbody tr').eq(2).children().eq(0).should('have.text', 'Reward Rate')
    cy.get('[class*="Block_block-table"] tbody tr').eq(2).children().eq(1).should('contain', '11.999961')
    cy.get('[class*="Block_block-table"] tbody tr').eq(3).children().eq(0).should('have.text', 'Proposer')
    cy.get('[class*="Block_block-table"] tbody tr').eq(3).children().eq(1).should('have.text', 'SCCXYUFQO54BDRMBOXV5GWGBARO2XA2H4TSNISLKK3DPFVGYCGQW7XNVA4')
    cy.get('[class*="Block_block-table"] tbody tr').eq(4).children().eq(0).should('have.text', 'Block Hash')
    cy.get('[class*="Block_block-table"] tbody tr').eq(4).children().eq(1).should('have.text', 'P3HVMIPSCDCCWLQYTTH46FDAOPVGPEBPCDH3BXWXGHJYA2IG3CNA')
    cy.get('[class*="Block_block-table"] tbody tr').eq(5).children().eq(0).should('have.text', 'Previous Block Hash')
    cy.get('[class*="Block_block-table"] tbody tr').eq(5).children().eq(1).should('have.text', 'PKDXKUHJNUPBHKHKEGVW2MA7LIGLMSUHS5GKOQ5YYMAO2OMK76EA')
    cy.get('[class*="Block_block-table"] tbody tr').eq(6).children().eq(0).should('have.text', 'Seed')
    cy.get('[class*="Block_block-table"] tbody tr').eq(6).children().eq(1).should('have.text', '9G7ZdsU9zCVlWW7uWBLxzWt2YTx0eXhLETbpNBea1Kc=')
  })

  it('shows page=1 on the url if not provided in url param', () => {
    cy.url().should('contain', '?page=1')
  })

  it('clicking txs list tx id navigates to tx page', () => {
    cy.get('[class*="Block_transactions-table"] tbody tr:first-child td:first-child a').as('txid')
    cy.get('@txid').should('have.text', 'OVVVKO...RONF6A')
    cy.get('@txid').click({force: true})
    cy.url().should('include', '/tx/')
  })

  it('clicking txs list From navigates to address page', () => {
    cy.get('[class*="Block_transactions-table"] tbody tr:first-child').children().eq(2).children().as('from')
    cy.get('@from').should('have.text', 'ZW3ISE...67W754')
    cy.get('@from').click({force: true})
    cy.url().should('include', '/address/')
  })

  it('clicking txs list To navigates to address page', () => {
    cy.get('[class*="Block_transactions-table"] tbody tr:first-child').children().eq(3).children().as('to')
    cy.get('@to').should('have.text', '6JUDNT...4NM6UE')
    cy.get('@to').click({force: true})
    cy.url().should('include', '/address/')
  })

  it('shows N/A as To for App Call type transaction', () => {
    cy.get('[class*="Block_transactions-table"] tbody tr').eq(2).children().eq(1).should('have.text', 'App Call')
    cy.get('[class*="Block_transactions-table"] tbody tr').eq(2).children().eq(3).should('have.text', 'N/A')
  })
})
