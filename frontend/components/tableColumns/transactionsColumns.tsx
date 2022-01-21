import React from "react";
import Link from "next/link";
import { Column, Row } from "react-table";
import TimeAgo from "timeago-react";
import AlgoIcon from "../algoicon";
import { TransactionResponse } from "../../types/apiResponseTypes";
import { IAsaMap } from "../../types/misc";
import {
  ellipseAddress,
  formatAsaAmountWithDecimal,
  formatNumber,
  getTxTypeName,
  integerFormatter,
  microAlgosToAlgos,
  removeSpace,
  TxType,
} from "../../utils/stringUtils";

export const transactionsColumns = (
  asaMap: IAsaMap,
  addr?: string,
  hideCols?: number[]
): Column[] => {
  let cols: Object[] | null[] = [
    {
      Header: "Tx ID",
      accessor: "id",
      Cell: ({ value }: { value: string }) => (
        <Link href={`/tx/${value}`}>{ellipseAddress(value)}</Link>
      ),
    },
    {
      Header: "Block",
      accessor: "confirmed-round",
      Cell: ({ value }: { value: number }) => {
        const _value = removeSpace(value.toString());
        return (
          <Link href={`/block/${_value}`}>
            {integerFormatter.format(Number(_value))}
          </Link>
        );
      },
    },
    {
      Header: "Type",
      accessor: "tx-type",
      Cell: ({ value }: { value: TxType }) => (
        <span className="type noselect">{getTxTypeName(value)}</span>
      ),
    },
    {
      Header: "From",
      accessor: "sender",
      Cell: ({ value }: { value: string }) =>
        value === addr ? (
          <span>{ellipseAddress(value)}</span>
        ) : (
          <Link href={`/address/${value}`}>{ellipseAddress(value)}</Link>
        ),
    },
    {
      Header: "To",
      accessor: "payment-transaction.receiver",
      Cell: ({ row }: { row: Row<TransactionResponse> }) => {
        const tx = row.original;
        const isAsaTransfer = tx["tx-type"] === TxType.AssetTransfer;
        const _value = isAsaTransfer
          ? tx["asset-transfer-transaction"].receiver
          : tx["payment-transaction"]
          ? tx["payment-transaction"].receiver
          : null;
        if (_value === addr) {
          return <span>{ellipseAddress(_value)}</span>;
        }
        return _value ? (
          <Link href={`/address/${_value}`}>{ellipseAddress(_value)}</Link>
        ) : (
          "N/A"
        );
      },
    },
    {
      Header: "Amount",
      accessor: "payment-transaction.amount",
      Cell: ({ row }: { row: Row<TransactionResponse> }) => {
        const tx = row.original;
        const _asaAmount =
          (tx["asset-transfer-transaction"] &&
            asaMap[tx["asset-transfer-transaction"]["asset-id"]] &&
            Number(
              formatAsaAmountWithDecimal(
                BigInt(tx["asset-transfer-transaction"].amount),
                asaMap[tx["asset-transfer-transaction"]["asset-id"]].decimals
              )
            )) ??
          0;
        const _asaUnit =
          tx["asset-transfer-transaction"] &&
          asaMap[tx["asset-transfer-transaction"]["asset-id"]] &&
          asaMap[tx["asset-transfer-transaction"]["asset-id"]].unitName;

        return (
          <span>
            {tx["tx-type"] === TxType.AssetTransfer ? (
              `${formatNumber(_asaAmount)} ${_asaUnit}`
            ) : tx["payment-transaction"] ? (
              <>
                <AlgoIcon />{" "}
                {formatNumber(
                  Number(microAlgosToAlgos(tx["payment-transaction"].amount))
                )}
              </>
            ) : (
              "N/A"
            )}
          </span>
        );
      },
    },
    {
      Header: "Fee",
      accessor: "fee",
      Cell: ({ value }: { value: number }) => (
        <span>
          <AlgoIcon /> {microAlgosToAlgos(value)}
        </span>
      ),
    },
    {
      Header: "Time",
      accessor: "round-time",
      Cell: ({ value }: { value: number }) => (
        <span>
          <TimeAgo datetime={new Date(value * 1000)} locale="en_short" />
        </span>
      ),
    },
  ];

  if (hideCols) {
    hideCols.forEach((colIndex) => {
      cols.splice(colIndex, 1, null);
    });
  }

  return cols.filter((col) => !!col) as Column[];
};
