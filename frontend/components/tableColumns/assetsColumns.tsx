import React from "react";
import Link from "next/link";
import {
  ellipseAddress,
  formatAsaAmountWithDecimal,
  formatNumber,
  integerFormatter,
  removeSpace,
} from "../../utils/stringUtils";
import { AccountOwnedAsset } from "../../types/apiResponseTypes";
import { Row } from "react-table";
import { IAsaMap } from "../../types/misc";

export const assetsColumns = (asaMap: IAsaMap, addr: string) => [
  {
    Header: "Asset ID",
    accessor: "asset-id",
    Cell: ({ value }: { value: number }) => (
      <Link href={`/asset/${value}`}>{removeSpace(value.toString())}</Link>
    ),
  },
  {
    Header: "Creator",
    accessor: "creator",
    Cell: ({ value }: { value: string }) =>
      value === addr ? (
        <span>{ellipseAddress(value)}</span>
      ) : (
        <Link href={`/address/${value}`}>{ellipseAddress(value)}</Link>
      ),
  },
  {
    Header: "Amount",
    accessor: "amount",
    Cell: ({ row }: { row: Row<AccountOwnedAsset> }) => {
      const asa = row.original;
      const _asaAmount =
        (asa.amount &&
          asaMap[asa["asset-id"]] &&
          Number(
            formatAsaAmountWithDecimal(
              BigInt(asa.amount),
              asaMap[asa["asset-id"]].decimals
            )
          )) ??
        0;
      const _asaUnit =
        asaMap[asa["asset-id"]] && asaMap[asa["asset-id"]].unitName;

      return (
        <span>
          {formatNumber(_asaAmount)} {_asaUnit}
        </span>
      );
    },
  },
  {
    Header: "Deleted",
    accessor: "deleted",
    Cell: ({ value }: { value: string }) => value.toString(),
  },
  {
    Header: "Is Frozen",
    accessor: "is-frozen",
    Cell: ({ value }: { value: string }) => value.toString(),
  },
  {
    Header: "Opted in at block",
    accessor: "opted-in-at-round",
    Cell: ({ value }: { value: string }) => (
      <Link href={`/block/${value}`}>
        {integerFormatter.format(Number(value))}
      </Link>
    ),
  },
];
