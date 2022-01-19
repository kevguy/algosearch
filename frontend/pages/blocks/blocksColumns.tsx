import Link from "next/link";
import { Row } from "react-table";
import TimeAgo from "timeago-react";
import AlgoIcon from "../../components/algoicon";
import { IBlockResponse, IBlockRewards } from "../../types/apiResponseTypes";
import {
  ellipseAddress,
  getTxTypeName,
  integerFormatter,
  microAlgosToAlgos,
  TxType,
} from "../../utils/stringUtils";
import styles from "./blocks.module.scss";

export const blocksColumns = [
  {
    Header: "Block",
    accessor: "round",
    Cell: ({ value }: { value: number }) => {
      const _value = value.toString().replace(" ", "");
      return (
        <Link href={`/block/${_value}`}>{integerFormatter.format(value)}</Link>
      );
    },
  },
  {
    Header: "Proposed by",
    accessor: "proposer",
    Cell: ({ value }: { value: string }) => (
      <Link href={`/address/${value}?page=1`}>{ellipseAddress(value)}</Link>
    ),
  },
  {
    Header: "# Tx",
    accessor: "transactions",
    Cell: ({ value }: { value: [] }) => {
      return <span>{value ? integerFormatter.format(value.length) : 0}</span>;
    },
  },
  {
    Header: "Transactions",
    accessor: "transaction-root",
    Cell: ({ row }: { row: Row<IBlockResponse> }) => {
      const txs = row.original.transactions;
      return txs && txs.length ? (
        <div className={styles["tx-type-wrapper"]}>
          {(Object.keys(TxType) as Array<keyof typeof TxType>).map((txType) => {
            const typeCount =
              txs &&
              txs.filter((tx) => tx["tx-type"] === TxType[txType]).length;
            return typeCount ? (
              <span key={txType}>
                {typeCount} {getTxTypeName(TxType[txType])}
                {typeCount > 1 && "s"}
              </span>
            ) : null;
          })}
        </div>
      ) : null;
    },
  },
  {
    Header: "Block Rewards",
    accessor: "rewards",
    Cell: ({ value }: { value: IBlockRewards }) => (
      <span>
        <AlgoIcon /> {microAlgosToAlgos(value["rewards-rate"])}
      </span>
    ),
  },
  {
    Header: "Time",
    accessor: "timestamp",
    Cell: ({ value }: { value: number }) => (
      <span>
        <TimeAgo datetime={new Date(value * 1000)} locale="en_short" />
      </span>
    ),
  },
];
