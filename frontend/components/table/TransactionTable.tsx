import axios from "axios";
import moment from "moment";
import Link from "next/link";
import React, { useEffect, useState } from "react";
import TimeAgo from "timeago-react";
import { TransactionResponse } from "../../pages/tx/[_txid]";
import { siteName } from "../../utils/constants";
import {
  ellipseAddress,
  getTxTypeName,
  microAlgosToAlgos,
  TxType,
} from "../../utils/stringUtils";
import AlgoIcon from "../algoicon";
import styles from "./TransactionTable.module.scss";

interface IAsaMap {
  [key: number]: string;
}

const TransactionTable = ({
  transactions,
}: {
  transactions: TransactionResponse[];
}) => {
  const [asaMap, setAsaMap] = useState<IAsaMap>([]);

  useEffect(() => {
    async function getAsas() {
      const dedupedAsaList = Array.from(
        new Set(
          transactions
            .filter((tx) => tx["tx-type"] === TxType.AssetTransfer)
            .map((tx) => tx["asset-transfer-transaction"]["asset-id"])
        )
      );
      const _asaList: string[] = await Promise.all(
        dedupedAsaList.map(
          async (asaId) =>
            await axios({
              method: "get",
              url: `${siteName}/v1/algod/assets/${asaId}`,
            })
              .then((response) => {
                console.log(
                  "asa unit name?",
                  response.data.params["unit-name"]
                );
                return response.data.params["unit-name"];
              })
              .catch((error) => {
                console.error("Error when retrieving Algorand ASA");
              })
        )
      );
      const _asaMap: IAsaMap = dedupedAsaList.reduce(
        (prev, asaId, index) => ({
          ...prev,
          [asaId]: _asaList[index],
        }),
        {}
      );
      if (_asaMap) {
        setAsaMap(_asaMap);
      }
      console.log("_asaMap: ", _asaMap);
    }
    getAsas();
  }, [transactions]);

  return (
    <div className={styles["transaction-table"]}>
      {transactions.map((tx: TransactionResponse, index: number) => {
        const _receiver = tx["payment-transaction"].receiver || tx.sender;
        let _asaUnit = asaMap[tx["asset-transfer-transaction"]["asset-id"]];

        return (
          <div key={tx.id} className={styles["transaction-row"]}>
            <div className={styles["transaction-subrow"]}>
              <span className={styles["transaction-id"]}>
                <Link href={`/transaction/${tx.id}`}>{tx.id}</Link>
              </span>
              <span className={styles.time}>
                <TimeAgo
                  datetime={new Date(moment.unix(tx["round-time"]).toDate())}
                  locale="en_short"
                />
              </span>
            </div>
            <div className={styles["transaction-subrow"]}>
              <div className={styles["relevant-accounts"]}>
                <span>
                  From:{" "}
                  <Link href={`/address/${tx.sender}`}>
                    {ellipseAddress(tx.sender)}
                  </Link>
                </span>
                <span>
                  To:{" "}
                  <Link href={`/address/${_receiver}`}>
                    {ellipseAddress(_receiver)}
                  </Link>
                </span>
              </div>
              <div className={styles["transaction-info"]}>
                <span>{getTxTypeName(tx["tx-type"])}</span>
                <span>
                  {tx["tx-type"] === TxType.AssetTransfer ? (
                    `${microAlgosToAlgos(
                      tx["asset-transfer-transaction"].amount
                    )} ${_asaUnit}`
                  ) : (
                    <>
                      <AlgoIcon />{" "}
                      {microAlgosToAlgos(tx["payment-transaction"].amount)}
                    </>
                  )}
                </span>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default TransactionTable;
