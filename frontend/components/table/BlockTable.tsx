import Link from "next/link";
import React from "react";
import TimeAgo from "timeago-react";
import { IBlockResponse } from "../../types/apiResponseTypes";
import {
  getTxTypeName,
  integerFormatter,
  TxType,
} from "../../utils/stringUtils";
import styles from "./BlockTable.module.scss";

const BlockTable = ({ blocks }: { blocks: IBlockResponse[] }) => {
  return (
    <div className={styles["block-table"]}>
      {blocks &&
        blocks.map((block: IBlockResponse, index: number) => {
          return (
            <div key={block._id} className={styles["block-row"]}>
              <div className={styles["block-subrow"]}>
                <span className={styles["block-id"]}>
                  <Link href={`/block/${block.round}`}>
                    {integerFormatter.format(block.round)}
                  </Link>
                </span>
                <span className={styles.proposer}>
                  Proposer:{" "}
                  {
                    <Link href={`/address/${block.proposer}`}>
                      {block.proposer}
                    </Link>
                  }
                </span>
                <span className={styles.time}>
                  <TimeAgo
                    datetime={new Date(block.timestamp * 1000)}
                    locale="en_short"
                  />
                </span>
              </div>
              <div className={styles["block-subrow"]}>
                <div className={styles["block-txs"]}>
                  <div className={styles["total-tx-label-wrapper"]}>
                    <span className={styles["total-tx-label"]}>
                      {block.transactions ? block.transactions.length : 0} Tx
                      {block.transactions && block.transactions.length > 1
                        ? "s"
                        : ""}
                    </span>
                  </div>
                  {block && block.transactions && block.transactions.length ? (
                    <div className={styles["txs-wrapper"]}>
                      {(Object.keys(TxType) as Array<keyof typeof TxType>).map(
                        (txType) => {
                          const typeCount =
                            block.transactions &&
                            block.transactions.filter(
                              (tx) => tx["tx-type"] === TxType[txType]
                            ).length;
                          return typeCount ? (
                            <span>
                              {typeCount} {getTxTypeName(TxType[txType])}
                              {typeCount > 1 && "s"}
                            </span>
                          ) : null;
                        }
                      )}
                    </div>
                  ) : null}
                </div>
              </div>
            </div>
          );
        })}
    </div>
  );
};

export default BlockTable;
