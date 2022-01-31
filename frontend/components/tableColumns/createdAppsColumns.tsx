import React from "react";
import Link from "next/link";
import { integerFormatter } from "../../utils/stringUtils";
import { StateSchema } from "../../types/apiResponseTypes";

export const createdAppsColumns = [
  {
    Header: "App ID",
    accessor: "id",
    Cell: ({ value }: { value: string }) => value,
  },
  {
    Header: "Created at Block",
    accessor: "created-at-round",
    Cell: ({ value }: { value: number }) => {
      const _value = value.toString().replace(" ", "");
      return (
        <Link href={`/block/${_value}`}>{integerFormatter.format(value)}</Link>
      );
    },
  },
  {
    Header: "Global State Schema",
    accessor: "params.global-state-schema",
    Cell: ({ value }: { value: StateSchema }) => {
      return (
        <div className="vertical-layout">
          <span># byte-slice: {value["num-byte-slice"]}</span>{" "}
          <span># uint: {value["num-uint"]}</span>
        </div>
      );
    },
  },
  {
    Header: "Local State Schema",
    accessor: "params.local-state-schema",
    Cell: ({ value }: { value: StateSchema }) => {
      return (
        <div className="vertical-layout">
          <span># byte-slice: {value["num-byte-slice"]}</span>{" "}
          <span># uint: {value["num-uint"]}</span>
        </div>
      );
    },
  },
  {
    Header: "Deleted",
    accessor: "deleted",
    Cell: ({ value }: { value: string }) => value.toString(),
  },
];
