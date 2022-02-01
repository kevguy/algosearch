import React from "react";
import Link from "next/link";
import {
  ellipseAddress,
  formatAsaAmountWithDecimal,
  formatNumber,
  isZeroAddress,
} from "../../utils/stringUtils";
import { IAsaResponse } from "../../types/apiResponseTypes";
import { Row } from "react-table";

export const createdAssetsColumns = (addr: string) => [
  {
    Header: "Asset ID",
    accessor: "index",
    Cell: ({ value }: { value: string }) => value,
  },
  {
    Header: "Asset Name",
    accessor: "params.name",
    Cell: ({ row }: { row: Row<IAsaResponse> }) => {
      const asa = row.original;

      return asa.params.url ? (
        <a href={asa.params.url} target="_blank" rel="noopener noreferrer">
          {asa.params.name}
        </a>
      ) : (
        asa.params.name
      );
    },
  },
  {
    Header: "Manager",
    accessor: "params.manager",
    Cell: ({ value }: { value: string }) =>
      !isZeroAddress(value) ? (
        value === addr ? (
          <span>{ellipseAddress(value)}</span>
        ) : (
          <Link href={`/address/${value}`}>{ellipseAddress(value)}</Link>
        )
      ) : (
        "N/A"
      ),
  },
  {
    Header: "Reserve",
    accessor: "params.reserve",
    Cell: ({ value }: { value: string }) =>
      !isZeroAddress(value) ? (
        value === addr ? (
          <span>{ellipseAddress(value)}</span>
        ) : (
          <Link href={`/address/${value}`}>{ellipseAddress(value)}</Link>
        )
      ) : (
        "N/A"
      ),
  },
  {
    Header: "Freeze",
    accessor: "params.freeze",
    Cell: ({ value }: { value: string }) =>
      !isZeroAddress(value) ? (
        value === addr ? (
          <span>{ellipseAddress(value)}</span>
        ) : (
          <Link href={`/address/${value}`}>{ellipseAddress(value)}</Link>
        )
      ) : (
        "N/A"
      ),
  },
  {
    Header: "Clawback",
    accessor: "params.clawback",
    Cell: ({ value }: { value: string }) =>
      !isZeroAddress(value) ? (
        value === addr ? (
          <span>{ellipseAddress(value)}</span>
        ) : (
          <Link href={`/address/${value}`}>{ellipseAddress(value)}</Link>
        )
      ) : (
        "N/A"
      ),
  },
  {
    Header: "Decimals",
    accessor: "params.decimals",
    Cell: ({ value }: { value: string }) => value,
  },
  {
    Header: "Total",
    accessor: "params.total",
    Cell: ({ row }: { row: Row<IAsaResponse> }) => {
      const asa = row.original;
      return (
        <>
          {formatNumber(
            Number(
              formatAsaAmountWithDecimal(
                BigInt(asa.params.total),
                asa.params.decimals
              )
            )
          )}{" "}
          {asa.params["unit-name"]}
        </>
      );
    },
  },
];
