import React, { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import AlgoIcon from "../../components/algoicon";
import {
  checkBase64EqualsEmpty,
  ellipseAddress,
  getTxTypeName,
  integerFormatter,
  microAlgosToAlgos,
  prettyPrintTEAL,
  removeSpace,
  TxType,
} from "../../utils/stringUtils";
import styles from "./TransactionDetails.module.scss";
import blockStyles from "../block/Block.module.scss";
import algosdk from "algosdk";
import msgpack from "@ygoe/msgpack";
import { TransactionResponse } from "../../types/apiResponseTypes";
import { IAsaMap } from "../../types/misc";
import { apiGetASA, getLsigTEAL } from "../../utils/api";
import {
  TabPanelUnstyled,
  TabsListUnstyled,
  TabsUnstyled,
  TabUnstyled,
} from "@mui/material";
import hljs from "highlight.js";
import TransactionAdditionalInfo from "../../components/transaction/TransactionAdditionalInfo";
import ApplicationTransactionInfo from "../../components/transaction/ApplicationTransactionInfo";
import {
  getAmount,
  getCloseAmount,
  getInnerTxCloseTo,
  getInnerTxReceiver,
} from "../../components/transaction/TransactionContentComponents";
import {
  algodAddr,
  algodProtocol,
  algodToken,
  isLocal,
} from "../../utils/constants";
import { DryrunResponse } from "algosdk/dist/types/src/client/v2/algod/models/types";

const TransactionDetails = ({
  transaction,
}: {
  transaction: TransactionResponse;
}) => {
  const [msgpackNotes, setMsgpackNotes] = useState();
  const [txType, setTxType] = useState<TxType>();
  const [receiver, setReceiver] = useState<string>();
  const [asaMap, setAsaMap] = useState<IAsaMap>([]);
  const [decodedNotes, setDecodedNotes] = useState<bigint>();
  const [disassembledLogicSig, setDisassembledLogicSig] = useState<string>();
  const decodeWithMsgpack = useCallback(() => {
    try {
      let message = msgpack.deserialize(
        Buffer.from(transaction.note, "base64")
      );
      if (typeof message === "object") {
        message = JSON.stringify(message, undefined, 2);
      }
      setMsgpackNotes(message);
    } catch (err) {
      return null;
    }
  }, [transaction]);

  useEffect(() => {
    if (transaction) {
      if (
        isLocal &&
        algodToken &&
        algodProtocol &&
        algodAddr &&
        transaction.signature.logicsig.logic &&
        transaction.signature.logicsig.args
      ) {
        const logicSig = new algosdk.LogicSigAccount(
          Buffer.from(transaction.signature.logicsig.logic, "base64"),
          transaction.signature.logicsig.args.map((item) =>
            Buffer.from(item, "base64")
          )
        );
        getLsigTEAL(logicSig, transaction)
          .then((result: DryrunResponse) => {
            if (
              result &&
              result.txns &&
              result.txns[0] &&
              result.txns[0].disassembly
            ) {
              const disassembledResult = prettyPrintTEAL(
                result.txns[0].disassembly
              );
              setDisassembledLogicSig(disassembledResult);
            }
          })
          .catch((error) => {
            console.error("LogicSig disassembly error: ", error);
          });
      }
      setTxType(transaction["tx-type"]);
      setReceiver(
        transaction && transaction["tx-type"] === TxType.AssetTransfer
          ? transaction["asset-transfer-transaction"].receiver
          : transaction["payment-transaction"]
          ? transaction["payment-transaction"].receiver
          : "N/A"
      );
      apiGetASA([transaction]).then((result) => {
        setAsaMap(result);
      });
      if (
        transaction["tx-type"] === TxType.App &&
        transaction["inner-txns"] &&
        transaction["inner-txns"].length > 0
      ) {
        apiGetASA(transaction["inner-txns"]).then((result) => {
          setAsaMap(result);
        });
      }
      if (
        transaction.note &&
        Buffer.from(transaction.note, "base64").length <= 8
      ) {
        setDecodedNotes(
          algosdk.decodeUint64(
            Buffer.from(transaction.note, "base64"),
            "bigint"
          )
        );
      }
      decodeWithMsgpack();
    }
  }, [decodeWithMsgpack, transaction]);
  if (!transaction) {
    return null;
  }
  return (
    <div className={blockStyles["table-wrapper"]}>
      <div className={blockStyles["block-table"]}>
        <table cellSpacing="0">
          <thead>
            <tr>
              <th>Identifier</th>
              <th>Value</th>
            </tr>
          </thead>
          <tbody>
            {transaction.group && !checkBase64EqualsEmpty(transaction.group) && (
              <tr>
                <td>Group ID</td>
                <td>{transaction.group}</td>
              </tr>
            )}
            <tr>
              <td>ID</td>
              <td>{transaction.id}</td>
            </tr>
            <tr>
              <td>Block</td>
              <td>
                <Link
                  href={`/block/${removeSpace(
                    transaction["confirmed-round"].toString()
                  )}`}
                >
                  {integerFormatter.format(
                    Number(
                      removeSpace(transaction["confirmed-round"].toString())
                    )
                  )}
                </Link>
              </td>
            </tr>
            <tr>
              <td>Type</td>
              <td>
                <span className="type noselect">{getTxTypeName(txType!)}</span>
              </td>
            </tr>
            <tr>
              <td>Sender</td>
              <td>
                <Link href={`/address/${transaction.sender}`}>
                  {transaction.sender}
                </Link>
              </td>
            </tr>
            {txType !== TxType.App && (
              <tr>
                <td>Receiver</td>
                <td>
                  {receiver ? (
                    <Link href={`/address/${receiver}`}>{receiver}</Link>
                  ) : (
                    "N/A"
                  )}
                </td>
              </tr>
            )}
            {txType !== TxType.App && (
              <tr>
                <td>Amount</td>
                <td>
                  <div>{getAmount(txType, transaction, asaMap)}</div>
                </td>
              </tr>
            )}
            <tr>
              <td>Fee</td>
              <td>
                <div>
                  <AlgoIcon /> {microAlgosToAlgos(transaction.fee)}
                </div>
              </td>
            </tr>
            <tr>
              <td>Timestamp</td>
              <td>{new Date(transaction["round-time"] * 1000).toString()}</td>
            </tr>
            {transaction.note && transaction.note !== "" && (
              <tr>
                <td className={styles["valign-top-identifier"]}>Note</td>
                <td>
                  <div>
                    <TabsUnstyled defaultValue={0}>
                      <TabsListUnstyled className={styles.tabs}>
                        <TabUnstyled>Base64</TabUnstyled>
                        <TabUnstyled>ASCII</TabUnstyled>
                        {decodedNotes && <TabUnstyled>UInt64</TabUnstyled>}
                        {msgpackNotes && <TabUnstyled>MessagePack</TabUnstyled>}
                      </TabsListUnstyled>
                      <TabPanelUnstyled value={0}>
                        <div className={styles.notes}>{transaction.note}</div>
                      </TabPanelUnstyled>
                      <TabPanelUnstyled value={1}>
                        <div className={styles.notes}>
                          {atob(transaction.note)}
                        </div>
                      </TabPanelUnstyled>
                      {decodedNotes && (
                        <TabPanelUnstyled value={2}>
                          <div className={styles["notes-row"]}>
                            <div>
                              <h5>Hexadecimal</h5>
                              <span>{decodedNotes!.toString(16)}</span>
                            </div>
                            <div>
                              <h5>Decimal</h5>
                              <span>{decodedNotes!.toString()}</span>
                            </div>
                          </div>
                        </TabPanelUnstyled>
                      )}
                      {msgpackNotes && (
                        <TabPanelUnstyled
                          value={!!decodedNotes ? 3 : 2}
                          className={styles.notes}
                        >
                          <pre>{msgpackNotes}</pre>
                        </TabPanelUnstyled>
                      )}
                    </TabsUnstyled>
                  </div>
                </td>
              </tr>
            )}
            {transaction.signature.logicsig.logic && disassembledLogicSig && (
              <tr>
                <td className={styles["valign-top-identifier"]}>LogicSig</td>
                <td>
                  <pre className={`${styles["teal-box"]} hljs`}>
                    <code
                      className="language-lua"
                      dangerouslySetInnerHTML={{
                        __html: hljs.highlight(disassembledLogicSig, {
                          language: "lua",
                        }).value,
                      }}
                    ></code>
                  </pre>
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
      {transaction["inner-txns"] && (
        <div>
          <h4>Inner Transactions</h4>
          <div
            className={`${blockStyles["block-table"]} ${styles["inner-txs-table"]}`}
          >
            <table cellSpacing="0">
              <thead>
                <tr>
                  <th>Type</th>
                  <th>Sender</th>
                  <th>Receiver</th>
                  <th>Amount</th>
                  <th>Close To</th>
                  <th>Close Amount</th>
                  <th>Fee</th>
                </tr>
              </thead>
              <tbody>
                {transaction["inner-txns"].map((innerTx, index) => (
                  <tr key={index}>
                    <td className={styles["normal-text"]}>
                      <h4 className="mobile-only">Type</h4>
                      {getTxTypeName(innerTx["tx-type"])}
                    </td>
                    <td>
                      <h4 className="mobile-only">Sender</h4>
                      <Link href={`/address/${innerTx.sender}`}>
                        {ellipseAddress(innerTx.sender)}
                      </Link>
                    </td>
                    <td>
                      <h4 className="mobile-only">Receiver</h4>
                      {getInnerTxReceiver(innerTx)}
                    </td>
                    <td>
                      <h4 className="mobile-only">Amount</h4>
                      {getAmount(innerTx["tx-type"], innerTx, asaMap)}
                    </td>
                    <td>
                      <h4 className="mobile-only">Close To</h4>
                      {getInnerTxCloseTo(innerTx)}
                    </td>
                    <td>
                      <h4 className="mobile-only">Close Amount</h4>
                      {getCloseAmount(innerTx["tx-type"], innerTx, asaMap)}
                    </td>
                    <td>
                      <h4 className="mobile-only">Fee</h4>
                      <AlgoIcon /> {microAlgosToAlgos(innerTx.fee)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
      {txType === TxType.App && (
        <ApplicationTransactionInfo transaction={transaction} />
      )}
      <TransactionAdditionalInfo transaction={transaction} />
    </div>
  );
};

export default TransactionDetails;
