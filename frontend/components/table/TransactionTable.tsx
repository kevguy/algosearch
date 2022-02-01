import Link from "next/link";
import React, { useEffect, useState } from "react";
import TimeAgo from "timeago-react";
import { TransactionResponse } from "../../types/apiResponseTypes";
import { IASAInfo, IAsaMap } from "../../types/misc";
import { apiGetASA } from "../../utils/api";
import {
  ellipseAddress,
  formatAsaAmountWithDecimal,
  formatNumber,
  getTxTypeName,
  microAlgosToAlgos,
  TxType,
} from "../../utils/stringUtils";
import AlgoIcon from "../algoicon";
import styles from "./TransactionTable.module.scss";

const TransactionTable = ({
  transactions,
}: {
  transactions?: TransactionResponse[];
}) => {
  const [asaMap, setAsaMap] = useState<IAsaMap>([]);

  useEffect(() => {
    if (!transactions) return;
    apiGetASA(transactions).then((result) => {
      setAsaMap(result);
    });
  }, [transactions]);

  if (!transactions) {
    return <></>;
  }

  return (
    <div className={styles["transaction-table"]}>
      {transactions &&
        transactions.map((tx: TransactionResponse, index: number) => {
          const _receiver =
            tx["payment-transaction"] && tx["payment-transaction"].receiver;
          let _asaInfo: IASAInfo =
            tx["asset-transfer-transaction"] &&
            asaMap[tx["asset-transfer-transaction"]["asset-id"]];

          return (
            <div key={tx.id} className={styles["transaction-row"]}>
              <div className={styles["transaction-subrow"]}>
                <div className={styles["tx-type-label-wrapper"]}>
                  <span className={styles["tx-type-label"]}>
                    {getTxTypeName(tx["tx-type"] as TxType)}
                  </span>
                </div>
                <span className={styles["transaction-id"]}>
                  <Link href={`/tx/${tx.id}`}>{tx.id}</Link>
                </span>
                <span className={styles.time}>
                  <TimeAgo
                    datetime={new Date(tx["round-time"] * 1000)}
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
                    {_receiver ? (
                      <Link href={`/address/${_receiver}`}>
                        {ellipseAddress(_receiver)}
                      </Link>
                    ) : (
                      "N/A"
                    )}
                  </span>
                </div>
                <div className={styles["transaction-info"]}>
                  <span>
                    {tx["tx-type"] === TxType.AssetTransfer &&
                    _asaInfo &&
                    tx["asset-transfer-transaction"] ? (
                      `${formatNumber(
                        Number(
                          formatAsaAmountWithDecimal(
                            BigInt(tx["asset-transfer-transaction"].amount),
                            _asaInfo.decimals
                          )
                        )
                      )} ${_asaInfo.unitName}`
                    ) : tx["payment-transaction"] ? (
                      <>
                        <AlgoIcon width={12} height={12} />{" "}
                        {formatNumber(
                          Number(
                            microAlgosToAlgos(tx["payment-transaction"].amount)
                          )
                        )}
                      </>
                    ) : (
                      "N/A"
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
