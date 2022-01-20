import React, { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import AlgoIcon from "../../components/algoicon";
import {
  formatAsaAmountWithDecimal,
  formatNumber,
  getTxTypeName,
  integerFormatter,
  microAlgosToAlgos,
  removeSpace,
  TxType,
} from "../../utils/stringUtils";
import styles from "./TransactionDetails.module.scss";
import blockStyles from "../block/Block.module.scss";
import algosdk from "algosdk";
import msgpack from "@ygoe/msgpack";
import { TransactionResponse } from "../../types/apiResponseTypes";
import { IAsaMap } from "../../types/misc";
import { apiGetASA } from "../../utils/api";
import {
  TabPanelUnstyled,
  TabsListUnstyled,
  TabsUnstyled,
  TabUnstyled,
} from "@mui/material";

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
      setTxType(transaction["tx-type"]);
      setReceiver(
        transaction && transaction["tx-type"] === TxType.AssetTransfer
          ? transaction["asset-transfer-transaction"].receiver
          : transaction["payment-transaction"].receiver
      );
      apiGetASA([transaction]).then((result) => {
        setAsaMap(result);
      });
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
            <tr>
              <td>Amount</td>
              <td>
                <div>
                  {txType === TxType.AssetTransfer ? (
                    <>
                      {transaction["asset-transfer-transaction"] &&
                        asaMap[
                          transaction["asset-transfer-transaction"]["asset-id"]
                        ] &&
                        formatNumber(
                          Number(
                            formatAsaAmountWithDecimal(
                              BigInt(
                                transaction["asset-transfer-transaction"].amount
                              ),
                              asaMap[
                                transaction["asset-transfer-transaction"][
                                  "asset-id"
                                ]
                              ].decimals
                            ) ?? 0
                          )
                        )}{" "}
                      {transaction["asset-transfer-transaction"] &&
                        asaMap[
                          transaction["asset-transfer-transaction"]["asset-id"]
                        ] && (
                          <Link
                            href={`/asset/${transaction["asset-transfer-transaction"]["asset-id"]}`}
                          >
                            {
                              asaMap[
                                transaction["asset-transfer-transaction"][
                                  "asset-id"
                                ]
                              ].unitName
                            }
                          </Link>
                        )}
                    </>
                  ) : (
                    <>
                      <AlgoIcon />{" "}
                      {formatNumber(
                        Number(
                          microAlgosToAlgos(
                            transaction["payment-transaction"].amount
                          )
                        )
                      )}
                    </>
                  )}
                </div>
              </td>
            </tr>
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
            <tr>
              <td className={styles["valign-top-identifier"]}>Note</td>
              <td>
                {transaction.note && transaction.note !== "" && (
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
                )}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div>
        <h4>Additional Information</h4>
        <div className={blockStyles["block-table"]}>
          <table cellSpacing="0">
            <thead>
              <tr>
                <th>Identifier</th>
                <th>Value</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>First Round</td>
                <td>
                  <Link
                    href={`/block/${removeSpace(
                      transaction["first-valid"].toString()
                    )}`}
                  >
                    {integerFormatter.format(
                      Number(removeSpace(transaction["first-valid"].toString()))
                    )}
                  </Link>
                </td>
              </tr>
              <tr>
                <td>Last Round</td>
                <td>
                  <Link
                    href={`/block/${removeSpace(
                      transaction["last-valid"].toString()
                    )}`}
                  >
                    {integerFormatter.format(
                      Number(removeSpace(transaction["last-valid"].toString()))
                    )}
                  </Link>
                </td>
              </tr>
              <tr>
                <td>Sender Rewards</td>
                <td>
                  <div>
                    <AlgoIcon />{" "}
                    {microAlgosToAlgos(transaction["sender-rewards"] || 0)}
                  </div>
                </td>
              </tr>
              <tr>
                <td>Receiver Rewards</td>
                <td>
                  <div>
                    <AlgoIcon />{" "}
                    {microAlgosToAlgos(transaction["receiver-rewards"] || 0)}
                  </div>
                </td>
              </tr>
              {Object.keys(transaction.signature.multisig).length > 0 && (
                <tr>
                  <td className={styles["valign-top-identifier"]}>Multisig</td>
                  <td className={styles["multisig-details"]}>
                    <div>Version {transaction.signature.multisig.version}</div>
                    <div>
                      Threshold: {transaction.signature.multisig.threshold}{" "}
                      signature
                      {transaction.signature.multisig.threshold! > 1 && "s"}
                    </div>
                    <h4>Subsignatures</h4>
                    {transaction.signature.multisig.subsignature?.map((sig) => {
                      const _addr = algosdk.encodeAddress(
                        Buffer.from(sig["public-key"], "base64")
                      );
                      return <Link href={`/address/${_addr}`}>{_addr}</Link>;
                    })}
                  </td>
                </tr>
              )}
              <tr>
                <td>Genesis ID</td>
                <td>{transaction["genesis-id"]}</td>
              </tr>
              <tr>
                <td>Genesis Hash</td>
                <td>{transaction["genesis-hash"]}</td>
              </tr>
              {transaction["tx-type"] === TxType.AssetConfig && (
                <>
                  <tr>
                    <td>Asset Name</td>
                    <td>
                      {transaction["asset-config-transaction"].params.url ? (
                        <a
                          href={
                            transaction["asset-config-transaction"].params.url
                          }
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          {transaction["asset-config-transaction"].params.name}
                        </a>
                      ) : (
                        transaction["asset-config-transaction"].params.name
                      )}
                    </td>
                  </tr>
                  <tr>
                    <td>Manager</td>
                    <td>
                      {transaction["asset-config-transaction"].params
                        .manager ? (
                        <Link
                          href={`/address/${transaction["asset-config-transaction"].params.manager}`}
                        >
                          {
                            transaction["asset-config-transaction"].params
                              .manager
                          }
                        </Link>
                      ) : (
                        "N/A"
                      )}
                    </td>
                  </tr>
                  <tr>
                    <td>Reserve</td>
                    <td>
                      {transaction["asset-config-transaction"].params
                        .reserve ? (
                        <Link
                          href={`/address/${transaction["asset-config-transaction"].params.reserve}`}
                        >
                          {
                            transaction["asset-config-transaction"].params
                              .reserve
                          }
                        </Link>
                      ) : (
                        "N/A"
                      )}
                    </td>
                  </tr>
                  <tr>
                    <td>Freeze</td>
                    <td>
                      {transaction["asset-config-transaction"].params.freeze ? (
                        <Link
                          href={`/address/${transaction["asset-config-transaction"].params.freeze}`}
                        >
                          {
                            transaction["asset-config-transaction"].params
                              .freeze
                          }
                        </Link>
                      ) : (
                        "N/A"
                      )}
                    </td>
                  </tr>
                  <tr>
                    <td>Clawback</td>
                    <td>
                      {transaction["asset-config-transaction"].params
                        .clawback ? (
                        <Link
                          href={`/address/${transaction["asset-config-transaction"].params.clawback}`}
                        >
                          {
                            transaction["asset-config-transaction"].params
                              .clawback
                          }
                        </Link>
                      ) : (
                        "N/A"
                      )}
                    </td>
                  </tr>
                  <tr>
                    <td>Decimals</td>
                    <td>
                      {transaction["asset-config-transaction"].params.decimals}
                    </td>
                  </tr>
                  <tr>
                    <td>Total</td>
                    <td>
                      {formatNumber(
                        Number(
                          formatAsaAmountWithDecimal(
                            BigInt(
                              transaction["asset-config-transaction"].params
                                .total
                            ),
                            transaction["asset-config-transaction"].params
                              .decimals
                          )
                        )
                      )}{" "}
                      {
                        transaction["asset-config-transaction"].params[
                          "unit-name"
                        ]
                      }
                    </td>
                  </tr>
                  <tr>
                    <td>Metadata Hash</td>
                    <td>
                      {
                        transaction["asset-config-transaction"].params[
                          "metadata-hash"
                        ]
                      }
                    </td>
                  </tr>
                </>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default TransactionDetails;
