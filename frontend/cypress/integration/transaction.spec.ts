import { backend_url, commonIntercepts, interceptDryrunEndpoint, stubWebSocketToInSync } from '../support/utils';
const txAcfgFixture = '../fixtures/tx/tx_acfg.json';
const txAfrzFixture = '../fixtures/tx/tx_afrz.json';
const txApprovalFixture = '../fixtures/tx/tx_appl_approval.json';
const txInnerTxsFixture = '../fixtures/tx/tx_appl_inner_txs_2.json';
const txAxferFixture = '../fixtures/tx/tx_axfer.json';
const txKeyregFixture = '../fixtures/tx/tx_keyreg.json';
const txLsigFixture = '../fixtures/tx/tx_pay_logicSig.json';
const txMsigFixture = '../fixtures/tx/tx_pay_multi.json';

describe('Transaction Page', () => {
  beforeEach(() => {
    commonIntercepts();

    stubWebSocketToInSync();
  });

  /* General, with Payment Tx with Single Sig */
  it('displays tx info correctly', () => {
    cy.visit('/tx/NTIU26TLJ6XMMBV6YQJB6SUPG5FBKCMHG2EQ5R5AGJDQ7OXK7PKQ');
    cy.wait('@getTx');

    cy.get('[class*="Block_block-table"]').eq(0).find('tbody tr').as('mainTable');
    cy.get('@mainTable').eq(0).children().eq(0).should('have.text', 'ID');
    cy.get('@mainTable').eq(0).children().eq(1).should('have.text', 'NTIU26TLJ6XMMBV6YQJB6SUPG5FBKCMHG2EQ5R5AGJDQ7OXK7PKQ');
    cy.get('@mainTable').eq(1).children().eq(0).should('have.text', 'Block');
    cy.get('@mainTable').eq(1).children().eq(1).should('have.text', '6,611,812');
    cy.get('@mainTable').eq(2).children().eq(0).should('have.text', 'Type');
    cy.get('@mainTable').eq(2).children().eq(1).should('have.text', 'Payment');
    cy.get('@mainTable').eq(3).children().eq(0).should('have.text', 'Sender');
    cy.get('@mainTable').eq(3).children().eq(1).should('have.text', 'TKEVVXMEEBRZQG5PP6CPJTHNK3MVHPTRA5M3ATER5GIL36EAPUMX5MR37U');
    cy.get('@mainTable').eq(4).children().eq(0).should('have.text', 'Receiver');
    cy.get('@mainTable').eq(4).children().eq(1).should('have.text', 'SP745JJR4KPRQEXJZHVIEN736LYTL2T2DFMG3OIIFJBV66K73PHNMDCZVM');
    cy.get('@mainTable').eq(5).children().eq(0).should('have.text', 'Amount');
    cy.get('@mainTable').eq(5).children().eq(1).should('contain', '28,883.027883');
    cy.get('@mainTable').eq(6).children().eq(0).should('have.text', 'Fee');
    cy.get('@mainTable').eq(6).children().eq(1).should('contain', '0.001');
  });

  it('displays notes in 4 versions', () => {
    cy.visit('/tx/NTIU26TLJ6XMMBV6YQJB6SUPG5FBKCMHG2EQ5R5AGJDQ7OXK7PKQ');
    cy.wait('@getTx');

    cy.get('[class*="Block_block-table"]').eq(0).find('tbody tr').eq(8).children().eq(0).should('have.text', 'Note');
    cy.get('[class*="Block_block-table"]').eq(0).find('tbody tr').eq(8).children().eq(1).as('noteTab');
    cy.get('@noteTab').find('.TabUnstyled-root').eq(0).should('have.text', 'Base64');
    cy.get('@noteTab').find('.TabUnstyled-root').eq(1).should('have.text', 'ASCII');
    cy.get('@noteTab').find('.TabUnstyled-root').eq(2).should('have.text', 'UInt64');
    cy.get('@noteTab').find('.TabUnstyled-root').eq(3).should('have.text', 'MessagePack');
    cy.get('@noteTab').find('.TabPanelUnstyled-root').eq(0).should('have.text', 'FszNghkCXmo=');
    cy.get('@noteTab').find('.TabUnstyled-root').eq(1).click();
    cy.get('@noteTab').find('.TabPanelUnstyled-root').eq(1).should('have.text', '\u0016ÌÍ\u0019\u0002^j');
    cy.get('@noteTab').find('.TabUnstyled-root').eq(2).click();
    cy.get('@noteTab').find('.TabPanelUnstyled-root').eq(2).find('[class*="TransactionDetails_notes-row"] div:first-child h5').should('have.text', 'Hexadecimal');
    cy.get('@noteTab').find('.TabPanelUnstyled-root').eq(2).find('[class*="TransactionDetails_notes-row"] div:first-child span').should('have.text', '16cccd8219025e6a');
    cy.get('@noteTab').find('.TabPanelUnstyled-root').eq(2).find('[class*="TransactionDetails_notes-row"] div:last-child h5').should('have.text', 'Decimal');
    cy.get('@noteTab').find('.TabPanelUnstyled-root').eq(2).find('[class*="TransactionDetails_notes-row"] div:last-child span').should('have.text', '1642913922732416618');
    cy.get('@noteTab').find('.TabUnstyled-root').eq(3).click();
    cy.get('@noteTab').find('.TabPanelUnstyled-root').eq(3).should('have.text', '22');
  });

  it('clicking Block ID navigates to block page', () => {
    cy.visit('/tx/NTIU26TLJ6XMMBV6YQJB6SUPG5FBKCMHG2EQ5R5AGJDQ7OXK7PKQ');
    cy.wait('@getTx');
    cy.get('[class*="Block_block-table"] tbody tr').eq(1).children().eq(1).children('a').click();
    cy.url().should('include', '/block/');
  });

  it('clicking Sender navigates to address page', () => {
    cy.visit('/tx/NTIU26TLJ6XMMBV6YQJB6SUPG5FBKCMHG2EQ5R5AGJDQ7OXK7PKQ');
    cy.wait('@getTx');
    cy.get('[class*="Block_block-table"] tbody tr').eq(3).children().eq(1).children('a').click();
    cy.url().should('include', '/address/');
  });

  it('clicking Receiver navigates to address page', () => {
    cy.visit('/tx/NTIU26TLJ6XMMBV6YQJB6SUPG5FBKCMHG2EQ5R5AGJDQ7OXK7PKQ');
    cy.wait('@getTx');
    cy.get('[class*="Block_block-table"] tbody tr').eq(4).children().eq(1).children('a').click();
    cy.url().should('include', '/address/');
  });

  it('displays Additional Info correctly', () => {
    cy.visit('/tx/NTIU26TLJ6XMMBV6YQJB6SUPG5FBKCMHG2EQ5R5AGJDQ7OXK7PKQ');
    cy.wait('@getTx');
    cy.get('[class*="Block_block-table"]').eq(1).find('tbody tr').as('addInfoTable');
    cy.get('@addInfoTable').eq(0).children().eq(0).should('have.text', 'First Valid');
    cy.get('@addInfoTable').eq(0).children().eq(1).should('have.text', '6,611,809');
    cy.get('@addInfoTable').eq(1).children().eq(0).should('have.text', 'Last Valid');
    cy.get('@addInfoTable').eq(1).children().eq(1).should('have.text', '6,612,809');
    cy.get('@addInfoTable').eq(2).children().eq(0).should('have.text', 'Confirmed Round');
    cy.get('@addInfoTable').eq(2).children().eq(1).should('have.text', '6,611,812');
    cy.get('@addInfoTable').eq(3).children().eq(0).should('have.text', 'Sender Rewards');
    cy.get('@addInfoTable').eq(3).children().eq(1).should('contain', '0.028883');
    cy.get('@addInfoTable').eq(4).children().eq(0).should('have.text', 'Receiver Rewards');
    cy.get('@addInfoTable').eq(4).children().eq(1).should('contain', '0');
    cy.get('@addInfoTable').eq(5).children().eq(0).should('have.text', 'Genesis ID');
    cy.get('@addInfoTable').eq(5).children().eq(1).should('have.text', 'mainnet-v1.0');
    cy.get('@addInfoTable').eq(6).children().eq(0).should('have.text', 'Genesis Hash');
    cy.get('@addInfoTable').eq(6).children().eq(1).should('have.text', 'wGHE2Pwdvd7S12BL5FaOP20EGYesN73ktiC1qzkkit8=');
  });

  /* Asset Config Tx */
  it('displays Asset Config Tx info correctly', () => {
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/transactions/*`,
      },
      {
        fixture: txAcfgFixture,
      }
    ).as('getAcfgTx');
    cy.visit('/tx/RR6ACOE4TJ6BOPEWFNWLE22IGX4WIG32O7BCMMFOKYTK5V7JQW5Q');
    cy.wait('@getAcfgTx');

    cy.get('[class*="Block_block-table"]').eq(0).find('tbody tr').as('acfgMainTable');
    cy.get('@acfgMainTable').eq(0).children().eq(0).should('have.text', 'ID');
    cy.get('@acfgMainTable').eq(0).children().eq(1).should('have.text', 'RR6ACOE4TJ6BOPEWFNWLE22IGX4WIG32O7BCMMFOKYTK5V7JQW5Q');
    cy.get('@acfgMainTable').eq(1).children().eq(0).should('have.text', 'Block');
    cy.get('@acfgMainTable').eq(1).children().eq(1).should('have.text', '4,162,690');
    cy.get('@acfgMainTable').eq(2).children().eq(0).should('have.text', 'Type');
    cy.get('@acfgMainTable').eq(2).children().eq(1).should('have.text', 'ASA Config');
    cy.get('@acfgMainTable').eq(3).children().eq(0).should('have.text', 'Sender');
    cy.get('@acfgMainTable').eq(3).children().eq(1).should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY');
    cy.get('@acfgMainTable').eq(4).children().eq(0).should('have.text', 'Receiver');
    cy.get('@acfgMainTable').eq(4).children().eq(1).should('have.text', 'N/A');
    cy.get('@acfgMainTable').eq(5).children().eq(0).should('have.text', 'Amount');
    cy.get('@acfgMainTable').eq(5).children().eq(1).should('contain', '0');
    cy.get('@acfgMainTable').eq(6).children().eq(0).should('have.text', 'Fee');
    cy.get('@acfgMainTable').eq(6).children().eq(1).should('contain', '0.001');

    /* Asset Config Info */
    cy.get('[class*="Block_block-table"]').eq(1).children('table').children('tbody').children('tr').as('acfgTableRows');
    cy.get('@acfgTableRows').eq(0).children().eq(0).should('have.text', 'Asset Name');
    cy.get('@acfgTableRows').eq(0).children().eq(1).should('have.text', 'Asia Reserve Currency Coin');
    cy.get('@acfgTableRows').eq(1).children().eq(0).should('have.text', 'Creator');
    cy.get('@acfgTableRows').eq(1).children().eq(1).should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY');
    cy.get('@acfgTableRows').eq(2).children().eq(0).should('have.text', 'Manager');
    cy.get('@acfgTableRows').eq(2).children().eq(1).should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY');
    cy.get('@acfgTableRows').eq(3).children().eq(0).should('have.text', 'Reserve');
    cy.get('@acfgTableRows').eq(3).children().eq(1).should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY');
    cy.get('@acfgTableRows').eq(4).children().eq(0).should('have.text', 'Freeze');
    cy.get('@acfgTableRows').eq(4).children().eq(1).should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY');
    cy.get('@acfgTableRows').eq(5).children().eq(0).should('have.text', 'Clawback');
    cy.get('@acfgTableRows').eq(5).children().eq(1).should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY');
    cy.get('@acfgTableRows').eq(6).children().eq(0).should('have.text', 'Decimals');
    cy.get('@acfgTableRows').eq(6).children().eq(1).should('have.text', '6');
    cy.get('@acfgTableRows').eq(7).children().eq(0).should('have.text', 'Total');
    cy.get('@acfgTableRows').eq(7).children().eq(1).should('have.text', '88,616,203,378.51 ARCC');
    cy.get('@acfgTableRows').eq(8).children().eq(0).should('have.text', 'Metadata Hash');
    cy.get('@acfgTableRows').eq(8).children().eq(1).should('have.text', 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=');
  });

  /* Asset Freeze Tx */
  it('displays Asset Freeze Tx info correctly', () => {
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/transactions/*`,
      },
      {
        fixture: txAfrzFixture,
      }
    ).as('getAfrzTx');
    cy.visit('/tx/BNNBILQQXACV2WT67EJINEW2RFA6ZZ4GP32TOKTPBJR2KXAS6EHA');
    cy.wait('@getAfrzTx');

    cy.get('[class*="Block_block-table"]').eq(0).find('tbody tr').as('afrzMainTable');
    cy.get('@afrzMainTable').eq(0).children().eq(0).should('have.text', 'Group ID');
    cy.get('@afrzMainTable').eq(0).children().eq(1).should('have.text', 'kvtthMqWp5BvLEJUToPtmnBs2begnkEyJ0iLT7G6um0=');
    cy.get('@afrzMainTable').eq(1).children().eq(0).should('have.text', 'ID');
    cy.get('@afrzMainTable').eq(1).children().eq(1).should('have.text', 'BNNBILQQXACV2WT67EJINEW2RFA6ZZ4GP32TOKTPBJR2KXAS6EHA');
    cy.get('@afrzMainTable').eq(2).children().eq(0).should('have.text', 'Block');
    cy.get('@afrzMainTable').eq(2).children().eq(1).should('have.text', '15,855,297');
    cy.get('@afrzMainTable').eq(3).children().eq(0).should('have.text', 'Type');
    cy.get('@afrzMainTable').eq(3).children().eq(1).should('have.text', 'ASA Freeze');
    cy.get('@afrzMainTable').eq(4).children().eq(0).should('have.text', 'Sender');
    cy.get('@afrzMainTable').eq(4).children().eq(1).should('have.text', 'FRWFUUDU5NLQQ2CC26VEILW3L3BCI7U6VNQJMTGIV43C7QEJ7ACCE3V4GM');
    cy.get('@afrzMainTable').eq(5).children().eq(0).should('have.text', 'Receiver');
    cy.get('@afrzMainTable').eq(5).children().eq(1).should('have.text', 'N/A');
    cy.get('@afrzMainTable').eq(6).children().eq(0).should('have.text', 'Amount');
    cy.get('@afrzMainTable').eq(6).children().eq(1).should('contain', 'N/A');
    cy.get('@afrzMainTable').eq(7).children().eq(0).should('have.text', 'Fee');
    cy.get('@afrzMainTable').eq(7).children().eq(1).should('contain', '0.001');

    /* Asset Freeze Info */
    cy.get('[class*="Block_block-table"]').eq(1).children('table').children('tbody').children('tr').as('afrzTableRows');
    cy.get('@afrzTableRows').eq(0).children().eq(0).should('have.text', 'Asset ID');
    cy.get('@afrzTableRows').eq(0).children().eq(1).should('have.text', '238498122');
    cy.get('@afrzTableRows').eq(1).children().eq(0).should('have.text', 'Address');
    cy.get('@afrzTableRows').eq(1).children().eq(1).should('have.text', '4EHY6C3XAXIUHIUTK2FQIKBCQLWFWMF7LNL2U3GXMIBHLHCXHJZPPKERTQ');
    cy.get('@afrzTableRows').eq(2).children().eq(0).should('have.text', 'New Freeze Status');
    cy.get('@afrzTableRows').eq(2).children().eq(1).should('have.text', 'false');
  });

  /* Asset Transfer Tx */
  it('displays Asset Transfer tx info correctly', () => {
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/transactions/*`,
      },
      {
        fixture: txAxferFixture,
      }
    ).as('getAxferTx');
    cy.visit('/tx/5FKLBBCQPBGAT4QWD6D4LANJ7GIXD2AVPK2W4PJBW77Q4U3LM5KQ');
    cy.wait('@getAxferTx');

    cy.get('[class*="Block_block-table"]').eq(0).find('tbody tr').as('axferTable');
    cy.get('@axferTable').eq(0).children().eq(0).should('have.text', 'ID');
    cy.get('@axferTable').eq(0).children().eq(1).should('have.text', '5FKLBBCQPBGAT4QWD6D4LANJ7GIXD2AVPK2W4PJBW77Q4U3LM5KQ');
    cy.get('@axferTable').eq(1).children().eq(0).should('have.text', 'Block');
    cy.get('@axferTable').eq(1).children().eq(1).should('have.text', '4,259,852');
    cy.get('@axferTable').eq(2).children().eq(0).should('have.text', 'Type');
    cy.get('@axferTable').eq(2).children().eq(1).should('have.text', 'ASA Transfer');
    cy.get('@axferTable').eq(3).children().eq(0).should('have.text', 'Sender');
    cy.get('@axferTable').eq(3).children().eq(1).should('have.text', 'ARCC3TMGVD7KXY7GYTE7U5XXUJXFRD2SXLAWRV57XJ6HWHRR37GNGNMPSY');
    cy.get('@axferTable').eq(4).children().eq(0).should('have.text', 'Receiver');
    cy.get('@axferTable').eq(4).children().eq(1).should('have.text', 'DBOVYJH5JCRIJ2UZIMY4QRSWXHHBTOOGUYQS7654WSJXAXQ2TUMWFVB2G4');
    cy.get('@axferTable').eq(5).children().eq(0).should('have.text', 'Amount');
    cy.get('@axferTable').eq(5).children().eq(1).should('contain', '888,888,888 ARCC');
    cy.get('@axferTable').eq(6).children().eq(0).should('have.text', 'Fee');
    cy.get('@axferTable').eq(6).children().eq(1).should('contain', '0.001');
  });

  /* Asset KeyReg Tx */
  it('displays Asset KeyReg tx info correctly', () => {
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/transactions/*`,
      },
      {
        fixture: txKeyregFixture,
      }
    ).as('getKeyregTx');
    cy.visit('/tx/QZBJYPNMCA5AMT6BFFA7QHEFC4JPPRGAX5DIXIUFPKRGX5EQURCA');
    cy.wait('@getKeyregTx');

    cy.get('[class*="Block_block-table"]').eq(0).find('tbody tr').eq(2).children().eq(0).should('have.text', 'Type');
    cy.get('[class*="Block_block-table"]').eq(0).find('tbody tr').eq(2).children().eq(1).should('have.text', 'Key Reg');

    cy.get('[class*="Block_block-table"]').eq(1).find('tbody tr').as('keyregTable');
    cy.get('@keyregTable').eq(0).children().eq(0).should('have.text', 'Mark account as participating');
    cy.get('@keyregTable').eq(0).children().eq(1).should('have.text', 'true');
    cy.get('@keyregTable').eq(1).children().eq(0).should('have.text', 'Selection Participation Key');
    cy.get('@keyregTable').eq(1).children().eq(1).should('have.text', 'YYWHHH53DCQZXYDMACPMC2DDF3DTYEBUUIMHGMAX27KZQHRQ2DKQ');
    cy.get('@keyregTable').eq(2).children().eq(0).should('have.text', 'Vote Participation Key');
    cy.get('@keyregTable').eq(2).children().eq(1).should('have.text', 'JD6AC6S6I4PHHXM65CDPSM2DHXIWZBKF3OPSTUTIZFKXA5JDLZGQ');
    cy.get('@keyregTable').eq(3).children().eq(0).should('have.text', 'Vote Key Dilution');
    cy.get('@keyregTable').eq(3).children().eq(1).should('have.text', '10,000');
    cy.get('@keyregTable').eq(4).children().eq(0).should('have.text', 'Vote First Valid');
    cy.get('@keyregTable').eq(4).children().eq(1).should('have.text', '15,777,579');
    cy.get('@keyregTable').eq(5).children().eq(0).should('have.text', 'Vote Last Valid');
    cy.get('@keyregTable').eq(5).children().eq(1).should('contain', '30,000,000');
  });

  /* App Call Tx with Approval Program and Clear State Program */
  it('displays App Call Tx with Approval Program and Clear State Program info correctly', () => {
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/transactions/*`,
      },
      {
        fixture: txApprovalFixture,
      }
    ).as('getApprovalTx');
    cy.visit('/tx/NTIU26TLJ6XMMBV6YQJB6SUPG5FBKCMHG2EQ5R5AGJDQ7OXK7PKQ');
    interceptDryrunEndpoint();
    cy.wait('@getApprovalTx');
    cy.wait('@getDryrunResponse');

    cy.get('[class*="Block_block-table"]').eq(0).find('tbody tr').as('appCallMainTable');
    cy.get('@appCallMainTable').eq(0).children().eq(0).should('have.text', 'Group ID');
    cy.get('@appCallMainTable').eq(0).children().eq(1).should('have.text', 'TwOfuW94Nd0hXrKfpSXcwD80DE7v1copZN8NZfQM67M=');
    cy.get('@appCallMainTable').eq(1).children().eq(0).should('have.text', 'ID');
    cy.get('@appCallMainTable').eq(1).children().eq(1).should('have.text', '3FLI5PVGYIHV7RXMW5M4VFVIIRQOZFMFELYGFL6FCEQUNOEAVD7Q');
    cy.get('@appCallMainTable').eq(2).children().eq(0).should('have.text', 'Block');
    cy.get('@appCallMainTable').eq(2).children().eq(1).should('have.text', '19,193,933');
    cy.get('@appCallMainTable').eq(3).children().eq(0).should('have.text', 'Type');
    cy.get('@appCallMainTable').eq(3).children().eq(1).should('have.text', 'App Call');
    cy.get('@appCallMainTable').eq(4).children().eq(0).should('have.text', 'Sender');
    cy.get('@appCallMainTable').eq(4).children().eq(1).should('have.text', 'MIZ3P3RZEXYT4VHRRT6K5EBYQKP24SLHBPXMADKKDZ3VCLVXOOKUACN42E');
    cy.get('@appCallMainTable').eq(5).children().eq(0).should('have.text', 'Fee');
    cy.get('@appCallMainTable').eq(5).children().eq(1).should('contain', '0.001');

    /* App Call Tx Info */
    cy.get('[class*="Block_block-table"]').eq(1).children('table').children('tbody').children('tr').as('appTxInfoTableRows');

    cy.get('@appTxInfoTableRows').eq(0).children().eq(0).should('have.text', 'Application ID');
    cy.get('@appTxInfoTableRows').eq(0).children().eq(1).should('have.text', '0');
    cy.get('@appTxInfoTableRows').eq(1).children().eq(0).should('have.text', 'Accounts');
    cy.get('@appTxInfoTableRows').eq(1).children().eq(1).should('have.text', 'N/A');
    cy.get('@appTxInfoTableRows').eq(2).children().eq(0).should('have.text', 'Arguments');
    cy.get('@appTxInfoTableRows').eq(2).children().eq(1).find('[class*="TransactionDetails_inner-table"]').as('appArgsTable');
    cy.get('@appArgsTable').find('thead td').eq(0).should('have.text', 'Base64');
    cy.get('@appArgsTable').find('thead td').eq(1).should('have.text', 'ASCII');
    cy.get('@appArgsTable').find('thead td').eq(2).should('have.text', 'UInt64');
    cy.get('@appArgsTable').find('tbody td').eq(0).should('have.text', 'Y3JlYXRl');
    cy.get('@appArgsTable').find('tbody td').eq(1).should('have.text', 'create');
    cy.get('@appArgsTable').find('tbody td').eq(2).should('have.text', '109342978307173');
    cy.get('@appTxInfoTableRows').eq(3).children().eq(0).should('have.text', 'On Completion');
    cy.get('@appTxInfoTableRows').eq(3).children().eq(1).should('have.text', 'noop');
    cy.get('@appTxInfoTableRows').eq(4).children().eq(0).should('have.text', 'Created Application Index');
    cy.get('@appTxInfoTableRows').eq(4).children().eq(1).should('have.text', '62368684');
    cy.get('@appTxInfoTableRows').eq(5).children().eq(0).should('have.text', 'Global State Schema');
    cy.get('@appTxInfoTableRows').eq(5).children().eq(1).should('contain', 'Number of byte-slice: 0');
    cy.get('@appTxInfoTableRows').eq(5).children().eq(1).should('contain', 'Number of uint: 0');
    cy.get('@appTxInfoTableRows').eq(6).children().eq(0).should('have.text', 'Local State Schema');
    cy.get('@appTxInfoTableRows').eq(6).children().eq(1).should('contain', 'Number of byte-slice: 0');
    cy.get('@appTxInfoTableRows').eq(6).children().eq(1).should('contain', 'Number of uint: 16');
    cy.get('@appTxInfoTableRows').eq(7).children().eq(0).should('have.text', 'Approval Program');
    cy.get('@appTxInfoTableRows').eq(7).children().eq(1).as('approvalTab');
    cy.get('@approvalTab').find('.TabUnstyled-root').eq(0).should('have.text', 'TEAL');
    cy.get('@approvalTab').find('.TabUnstyled-root').eq(1).should('have.text', 'Base64');
    cy.get('@approvalTab').find('.TabUnstyled-root').eq(1).click();
    cy.get('@approvalTab')
      .find('.TabPanelUnstyled-root')
      .eq(1)
      .should(
        'have.text',
        'BCAHAAHoB+UHBf///////////wHAhD0mDQFvAWUBcAJhMQJhMgJsdARzd2FwBG1pbnQBdAJjMQJwMQJjMgJwMjEZgQQSMRkhBBIRMRmBAhIRQATxMRkjEjEbIhIQQATjNhoAgAZjcmVhdGUSQATUMRkjEjYaAIAJYm9vdHN0cmFwEhBAA/MzAhIzAggINTQiK2I1ZSI0ZXAARDUBIicEYjVmNGZAABEiYCJ4CTEBCDMACAk1AkIACCI0ZnAARDUCIicFYjVnKDRlFlA1byI0b2I1PSg0ZhZQNXAiNHBiNT4oNGcWUDVxIjRxYjU/IipiNUA0ATQ9CTVHNAI0Pgk1SDEAKVA0ZRZQNXkxAClQNGYWUDV6MQApUDRnFlA1ezYaAIAGcmVkZWVtEkAAWjYaAIAEZmVlcxJAABw2GgAnBhI2GgAnBxIRNhoAgARidXJuEhFAAG0ANGdJRDMCERJEMwISRDMCFDIJEkQ0PzMCEgk1PzRAMwISCTVAIio0QGYiNHE0P2YjQzMCFDMCBzMCECMSTTYcARJENDREIigzAhEWUEpiNDQJZiMxAClQMwIRFlBKYjQ0CUlBAANmI0NIaCNDMgciJwhiCUk1+kEARiInCWIiJwpiNPodTEAANx4hBSMeHzX7SEhIIicLYiInDGI0+h1MQAAdHiEFIx4fNfxISEgiJwk0+2YiJws0/GYiJwgyB2YzAxIzAwgINTU2HAExABNENGdBACIiNGdwAEQ1BiIcNAYJND8INQQ2GgAnBhJAASA0ZzMEERJENhoAJwcSQABVNhwBMwQAEkQzBBI0Rx00BCMdH0hITEhJNRA0NAk1yTMEEjRIHTQEIx0fSEhMSEk1ETQ1CTXKNBA0ERBENEc0EAk1UTRINBEJNVI0BDMEEgk1U0ICCjYcATMCABJENEc0NAg1UTRINDUINVI0BCISQAAuNDQ0BB00RyMdH0hITEg0NTQEHTRIIx0fSEhMSEoNTUk0BAg1UzMEEgk1y0IBvyInBTMEEUk1Z2YoNGcWUDVxIjRncABERDRnNGUTRDRnNGYTRDMEEiQISR018DQ0NDUdNfFKDEAACBJENPA08Q5EMwQSJAgjCEkdNfA0NDQ1HTXxSg1AAAgSRDTwNPENRCQ1PzQEMwQSJAgINVNCAU82HAEzAgASRDMCETRlEjMDETRmEhBJNWRAABkzAhE0ZhIzAxE0ZRIQRDRINRI0RzUTQgAINEc1EjRINRM2GgGAAmZpEkAAWjYaAYACZm8SRDQ1JAs0Eh00EzQ1CSUdH0hITEgjCEk1FSINNDU0EwwQRDQ0NBUJNGRBABM1yTRHNBUINVE0SDQ1CTVSQgBnNco0SDQVCDVSNEc0NQk1UUIAVDQ0STUVJQs0Ex00EiQLNDQlCx4fSEhMSEk1FCINNBQ0EwwQRDQUNDUJNGRBABM1yjRHNDQINVE0SDQUCTVSQgATNck0RzQUCTVRNEg0NAg1UkIAADQVIQQLNAQdgaCcATQSHR9ISExISTUqNAQINVNCADsiKzYaARdJNWVmIicENhoCF0k1ZmY0ZXEDRIABLVCABEFMR080ZkEABkg0ZnEDRFAzAiZJFYEPTFISQyIqNEA0KghmIjRxND80Kgg0ywhmIjRvND00yQhmIjRwND40yghmIoACczE0UWYigAJzMjRSZiInCjRSIQYdNFEjHR9ISExIZiInDDRRIQYdNFIjHR9ISExIZiKAA2lsdDRTZjTLQQAJIzR7SmI0ywhmNMlBAAkjNHlKYjTJCGY0ykEACSM0ekpiNMoIZiNDI0MiQw=='
      );
    cy.get('@appTxInfoTableRows').eq(8).children().eq(0).should('have.text', 'Clear State Program');
    cy.get('@appTxInfoTableRows').eq(8).children().eq(1).as('clearTab');
    cy.get('@clearTab').find('.TabUnstyled-root').eq(0).should('have.text', 'TEAL');
    cy.get('@clearTab').find('.TabUnstyled-root').eq(1).should('have.text', 'Base64');
    cy.get('@clearTab').find('.TabUnstyled-root').eq(1).click();
    cy.get('@clearTab').find('.TabPanelUnstyled-root').eq(1).should('have.text', 'BIEB');
  });

  /* App Call Tx with Inner Txs */
  it('displays App Call Tx with Inner Txs info correctly', () => {
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/transactions/*`,
      },
      {
        fixture: txInnerTxsFixture,
      }
    ).as('getInnersTx');
    cy.visit('/tx/Y5W6UWSSYEEOXJI2RHID66HYWIBQ4EVCCCZY5Q4LHK6XSC74BFNA');
    cy.wait('@getInnersTx');

    cy.get('[class*="Block_block-table"]').eq(0).find('tbody tr').as('appCallMainTable');
    cy.get('@appCallMainTable').eq(0).children().eq(0).should('have.text', 'ID');
    cy.get('@appCallMainTable').eq(0).children().eq(1).should('have.text', 'Y5W6UWSSYEEOXJI2RHID66HYWIBQ4EVCCCZY5Q4LHK6XSC74BFNA');
    cy.get('@appCallMainTable').eq(1).children().eq(0).should('have.text', 'Block');
    cy.get('@appCallMainTable').eq(1).children().eq(1).should('have.text', '13,348,578');
    cy.get('@appCallMainTable').eq(2).children().eq(0).should('have.text', 'Type');
    cy.get('@appCallMainTable').eq(2).children().eq(1).should('have.text', 'App Call');
    cy.get('@appCallMainTable').eq(3).children().eq(0).should('have.text', 'Sender');
    cy.get('@appCallMainTable').eq(3).children().eq(1).should('have.text', '4WMWVTICJG2XWSZBXSRALIULEUWKKERMT2NSY67JGOFRKVEW7FURVYC5OM');
    cy.get('@appCallMainTable').eq(4).children().eq(0).should('have.text', 'Fee');
    cy.get('@appCallMainTable').eq(4).children().eq(1).should('contain', '0.001');

    /* Inner Txs Info */
    cy.get('[class*="Block_block-table"]').eq(1).children('table').children('tbody').children('tr').as('innersTableRows');
    cy.get('@innersTableRows').eq(0).children().eq(0).should('contain', 'ASA Transfer');
    cy.get('@innersTableRows').eq(0).children().eq(1).should('contain', 'BHBJZQ...XWYV4Q');
    cy.get('@innersTableRows').eq(0).children().eq(2).should('contain', 'AAAAAA...Y5HFKQ');
    cy.get('@innersTableRows').eq(0).children().eq(3).should('contain', '0 ARCC');
    cy.get('@innersTableRows').eq(0).children().eq(4).should('contain', 'BQLHUS...TTNIW4');
    cy.get('@innersTableRows').eq(0).children().eq(5).should('contain', '0.0002 ARCC');
    cy.get('@innersTableRows').eq(0).children().eq(6).should('contain', '0.001');
    cy.get('@innersTableRows').eq(1).children().eq(0).should('contain', 'Payment');
    cy.get('@innersTableRows').eq(1).children().eq(1).should('contain', 'BHBJZQ...XWYV4Q');
    cy.get('@innersTableRows').eq(1).children().eq(2).should('contain', 'AAAAAA...Y5HFKQ');
    cy.get('@innersTableRows').eq(1).children().eq(3).should('contain', '0');
    cy.get('@innersTableRows').eq(1).children().eq(4).should('contain', '4WMWVT...VYC5OM');
    cy.get('@innersTableRows').eq(1).children().eq(5).should('contain', '1.2');
    cy.get('@innersTableRows').eq(1).children().eq(6).should('contain', '0.001');

    /* App Call Tx Info */
    cy.get('[class*="Block_block-table"]').eq(2).children('table').children('tbody').children('tr').as('appTxInfoTableRows');

    cy.get('@appTxInfoTableRows').eq(0).children().eq(0).should('have.text', 'Application ID');
    cy.get('@appTxInfoTableRows').eq(0).children().eq(1).should('have.text', '400404272');
    cy.get('@appTxInfoTableRows').eq(1).children().eq(0).should('have.text', 'Accounts');
    cy.get('@appTxInfoTableRows').eq(1).children().eq(1).should('contain', '4WMWVTICJG2XWSZBXSRALIULEUWKKERMT2NSY67JGOFRKVEW7FURVYC5OM');
    cy.get('@appTxInfoTableRows').eq(1).children().eq(1).should('contain', 'BQLHUSSAXDPOAXLEYEHS332FIT2RKDQ33RAL45JVTXQYMKNZ3CB3TTNIW4');
    cy.get('@appTxInfoTableRows').eq(2).children().eq(0).should('have.text', 'Foreign Assets');
    cy.get('@appTxInfoTableRows').eq(2).children().eq(1).should('have.text', '27165954');
    cy.get('@appTxInfoTableRows').eq(3).children().eq(0).should('have.text', 'On Completion');
    cy.get('@appTxInfoTableRows').eq(3).children().eq(1).should('have.text', 'delete');
    cy.get('@appTxInfoTableRows').eq(4).children().eq(0).should('have.text', 'Global State Schema');
    cy.get('@appTxInfoTableRows').eq(4).children().eq(1).should('contain', 'Number of byte-slice: 0');
    cy.get('@appTxInfoTableRows').eq(4).children().eq(1).should('contain', 'Number of uint: 0');
    cy.get('@appTxInfoTableRows').eq(5).children().eq(0).should('have.text', 'Local State Schema');
    cy.get('@appTxInfoTableRows').eq(5).children().eq(1).should('contain', 'Number of byte-slice: 0');
    cy.get('@appTxInfoTableRows').eq(5).children().eq(1).should('contain', 'Number of uint: 0');
  });

  /* Payment Tx with LogicSig */
  it('displays LogicSig info correctly', () => {
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/transactions/*`,
      },
      {
        fixture: txLsigFixture,
      }
    ).as('getLsigTx');
    cy.visit('/tx/Z24TC75TMSLTN5X3WCKJGMW73YG2CLUU2A2EQG6E6MYNBZNO7KJA');
    interceptDryrunEndpoint();
    cy.wait('@getLsigTx');
    cy.wait('@getDryrunResponse');

    /* LogicSig Info */
    cy.get('[class*="Block_block-table"]').eq(0).children('table').children('tbody').children('tr').as('lsigTableRows');
    cy.get('@lsigTableRows').eq(7).children().eq(0).should('have.text', 'Close Amount');
    cy.get('@lsigTableRows').eq(7).children().eq(1).should('contain', '0.414609');
    cy.get('@lsigTableRows').eq(8).children().eq(0).should('have.text', 'Close Remainder To');
    cy.get('@lsigTableRows').eq(8).children().eq(1).should('have.text', '75IHRONIM4QX2RXSF5JIHH5GXVPQSRLR4CFJC4JUJJJHZPAYJIYDPQUKNE');
    cy.get('@lsigTableRows').eq(11).children().eq(0).should('have.text', 'LogicSig');
    cy.get('@lsigTableRows').eq(11).children().eq(1).as('lsigTab');
    cy.get('@lsigTab').find('.TabUnstyled-root').eq(0).should('have.text', 'TEAL');
    cy.get('@lsigTab').find('.TabUnstyled-root').eq(1).should('have.text', 'Base64');
    cy.get('@lsigTab').find('.TabUnstyled-root').eq(1).click();
    cy.get('@lsigTab')
      .find('.TabPanelUnstyled-root')
      .eq(1)
      .should(
        'have.text',
        'BCAIAQAEAwa/xsEKAovu1hEmASD/UHi5qGchfUbyL1KDn6a9XwlFceCKkXE0SlJ8vBhKMDIEJA5EMQGB6AcORCM1CTQJOCAyAxJENAk4FTIDEkQ0CSIINQk0CTIEDED/4jIEIQYSMgQlEhEzABAiEhAzARAhBBIQMwAIgaDCHg8QMwEIIxIQMwAJMgMSEDMBCTIDEhAzABkjEhAzARkiEhAzARghBRIQMwAAKBIQMwEAMQASEDMABzEAEhA1ADIEIQYSNQE0AUAAIDMCECQSMwISIxIQMwIAKBIQMwIZIxIQIQczAhESEDUBNAA0ARBBAAIiQzIEJRIzABghBRIQMwAJMgMSEDMBCSgSEDMCCTIDEhAzAAAxABIQMwEAMQASEDMCACgSEDMAECEEEhAzARAiEhAzAhAiEhAzAAgjEhAzAQgjEhAzAggjEhAzABklEhAzARkjEhAzAhkjEhBBAAIiQzMBCTIDEkAAYDMAGSEGEjIEJRIQMwAQIQQSEDMBECISEDMCECQSEDMAGCEFEhAzAAAxABIQMwEAMQASEDMCADEAExAzAQczAgASEDMCFCgSEDMACTIDEhAzAQkoEhAzAgkyAxIQREIAlTMAGSMSMQkyAxIQMgQkEhAzABAhBBIQMwEQIhIQMwIQJBIQMwMQIhIQMwMIgdAPDxAzAwcxABIQMwMAMwEHEhAzAAAxABIQMwEAMQASEDMCADEAExAzAwAxABMQMwAYIQUSEDMBADEAEhAzAQczAgASEDMCFCgSEDMACTIDEhAzAQkyAxIQMwIJMgMSEDMDCTIDEhBEMwEIIg8zAhIiDxAhBzMCERIQRDMCEiIdNQI1ATMBCCIdNQQ1AzQBNAMNQAAPNAE0AxI0AjQEDxBAAAEAIkM='
      );
  });

  /* Payment Tx with MultiSig */
  it('displays MultiSig info correctly', () => {
    cy.intercept(
      {
        method: 'GET',
        url: `${backend_url}/v1/transactions/*`,
      },
      {
        fixture: txMsigFixture,
      }
    ).as('getMsigTx');
    cy.visit('/tx/7GCCBF4LIVFVADLWV3TQHZYH7BS3ROUISJZHUI675WGGP5YVGH4Q');
    cy.wait('@getMsigTx');

    /* MultiSig Info */
    cy.get('[class*="Block_block-table"]').eq(1).children('table').children('tbody').children('tr').as('msigTableRows');
    cy.get('@msigTableRows').eq(5).children().eq(0).should('have.text', 'Multisig');
    cy.get('@msigTableRows').eq(5).children().eq(1).should('contain', 'Version 1');
    cy.get('@msigTableRows').eq(5).children().eq(1).should('contain', 'Threshold: 1 signature');
    cy.get('@msigTableRows').eq(5).children().eq(1).should('contain', 'P3MWU63PMOKDVWRFGZM2XX2CIQC4REUN277QNQWZZDN7Z24N6SW5VC5CBI');
    cy.get('@msigTableRows').eq(5).children().eq(1).should('contain', 'GSS25RXJLMKU7BJM46BJQPV7JWWRVAZMQG4YKTQMFTDYQEGGAONMDP2GW4');
  });
});
