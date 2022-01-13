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
import styles from "../block/Block.module.scss";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import Box from "@mui/material/Box";
import TabPanel from "../../components/tabPanel";
import algosdk from "algosdk";
import msgpack from "@ygoe/msgpack";
import { TransactionResponse } from "../../types/apiResponseTypes";
import { IAsaMap } from "../../types/misc";
import { apiGetASA } from "../../utils/api";

function a11yProps(index: number) {
  return {
    id: `simple-tab-${index}`,
    "aria-controls": `simple-tabpanel-${index}`,
  };
}

const TransactionDetails = ({
  transaction,
}: {
  transaction: TransactionResponse;
}) => {
  const [noteTab, setNoteTab] = useState(0);
  const [msgpackNotes, setMsgpackNotes] = useState();
  const [txType, setTxType] = useState<TxType>();
  const [receiver, setReceiver] = useState<string>();
  const [asaMap, setAsaMap] = useState<IAsaMap>([]);
  const [decodedNotes, setDecodedNotes] = useState<bigint>();
  const clickTabHandler = (event: React.SyntheticEvent, newValue: number) => {
    setNoteTab(newValue);
  };
  const decodeWithMsgpack = useCallback(() => {
    try {
      return msgpack.deserialize(Buffer.from(transaction.note, "base64"));
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
        Buffer.from(transaction.note, "base64").length < 8
      ) {
        setDecodedNotes(
          algosdk.decodeUint64(
            Buffer.from(transaction.note, "base64"),
            "bigint"
          )
        );
      }
      setMsgpackNotes(decodeWithMsgpack());
    }
  }, [decodeWithMsgpack, transaction]);
  if (!transaction) {
    return null;
  }
  return (
    <div className={styles["table-wrapper"]}>
      <div className={styles["block-table"]}>
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
              <td>Round</td>
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
                      {microAlgosToAlgos(
                        transaction["payment-transaction"].amount
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
              <td>Timestamp</td>
              <td>{new Date(transaction["round-time"] * 1000).toString()}</td>
            </tr>
            <tr>
              <td>Note</td>
              <td>
                {transaction.note && transaction.note !== "" && (
                  <div>
                    <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
                      <Tabs
                        value={noteTab}
                        onChange={clickTabHandler}
                        aria-label="Note in different encoding"
                      >
                        <Tab label="Base64" {...a11yProps(0)} />
                        <Tab label="ASCII" {...a11yProps(1)} />
                        {decodedNotes && (
                          <Tab label="Uint64" {...a11yProps(2)} />
                        )}
                        {msgpackNotes && (
                          <Tab label="MessagePack" {...a11yProps(3)} />
                        )}
                      </Tabs>
                    </Box>
                    <TabPanel value={noteTab} index={0}>
                      <div className={styles.notes}>{transaction.note}</div>
                    </TabPanel>
                    <TabPanel value={noteTab} index={1}>
                      <div className={styles.notes}>
                        {atob(transaction.note)}
                      </div>
                    </TabPanel>
                    {decodedNotes && (
                      <TabPanel value={noteTab} index={2}>
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
                      </TabPanel>
                    )}
                    {msgpackNotes && (
                      <TabPanel value={noteTab} index={decodedNotes ? 3 : 2}>
                        <div className={styles.notes}>{msgpackNotes}</div>
                      </TabPanel>
                    )}
                  </div>
                )}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div>
        <h4>Miscellaneous Details</h4>
        <div className={styles["block-table"]}>
          <table cellSpacing="0">
            <thead>
              <tr>
                <th>Identifier</th>
                <th>Value</th>
              </tr>
            </thead>
            <tbody>
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
