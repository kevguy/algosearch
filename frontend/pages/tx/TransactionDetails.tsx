import React, { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import AlgoIcon from "../../components/algoicon";
import {
  getTxTypeName,
  integerFormatter,
  microAlgosToAlgos,
  removeSpace,
  TxType,
} from "../../utils/stringUtils";
import "react-table-6/react-table.css";
import styles from "../block/Block.module.scss";
import moment from "moment";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import Box from "@mui/material/Box";
import TabPanel from "../../components/tabPanel";
import algosdk from "algosdk";
import msgpack from "@ygoe/msgpack";
import { TransactionResponse } from "../../types/apiResponseTypes";

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
  let txType: TxType | undefined;
  let receiver;
  let decodedNotes: bigint | undefined;
  const clickTabHandler = (event: React.SyntheticEvent, newValue: number) => {
    setNoteTab(newValue);
  };
  const decodeWithMsgpack = useCallback(() => {
    try {
      return msgpack.deserialize(Buffer.from(transaction.note, "base64"));
    } catch (err) {
      return null;
    }
  }, []);
  useEffect(() => {
    if (transaction) {
      txType = transaction["tx-type"];
      receiver =
        transaction && txType === TxType.AssetTransfer
          ? transaction["asset-transfer-transaction"].receiver
          : transaction["payment-transaction"].receiver;
      if (Buffer.from(transaction.note, "base64").length < 8) {
        decodedNotes = algosdk.decodeUint64(
          Buffer.from(transaction.note, "base64"),
          "bigint"
        );
      }
      setMsgpackNotes(decodeWithMsgpack());
    }
  }, [decodeWithMsgpack]);
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
                  <AlgoIcon />{" "}
                  {txType === TxType.AssetTransfer
                    ? transaction["asset-transfer-transaction"].amount // need to divide by decimal
                    : microAlgosToAlgos(
                        transaction["payment-transaction"].amount
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
              <td>First round</td>
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
              <td>Last round</td>
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
              <td>{moment.unix(transaction["round-time"]).format("LLLL")}</td>
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
                      {transaction.note}
                    </TabPanel>
                    <TabPanel value={noteTab} index={1}>
                      {atob(transaction.note)}
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
                        {msgpackNotes}
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
                <td>Sender rewards</td>
                <td>
                  <div>
                    <AlgoIcon />{" "}
                    {microAlgosToAlgos(transaction["sender-rewards"] || 0)}
                  </div>
                </td>
              </tr>
              <tr>
                <td>Receiver rewards</td>
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
                <td>Genesis hash</td>
                <td>{transaction["genesis-hash"]}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default TransactionDetails;
