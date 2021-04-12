import React from "react";
import MaUTable from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import { Column, useTable } from "react-table";
import styles from "./TransactionTable.module.scss";

const Table = ({
  columns,
  data,
  loading,
  getProps,
  className,
}: {
  columns: Column[];
  data: readonly {}[];
  loading?: boolean;
  getProps?: Function;
  className?: string;
}) => {
  // Use the state and functions returned from useTable to build your UI
  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } =
    useTable({
      columns,
      data,
    });

  // Render the UI for your table
  return (
    <>
      <MaUTable
        {...getTableProps()}
        className={`${styles["mui-table"]}${className ? " " + className : ""}`}
      >
        <TableHead>
          {headerGroups.map((headerGroup) => (
            <TableRow
              {...headerGroup.getHeaderGroupProps()}
              key={headerGroup.getHeaderGroupProps().key}
            >
              {headerGroup.headers.map((column) => (
                <TableCell
                  {...column.getHeaderProps()}
                  key={column.getHeaderProps().key}
                >
                  {column.render("Header")}
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableHead>
        <TableBody {...getTableBodyProps()}>
          {rows.map((row, i) => {
            prepareRow(row);
            return (
              <TableRow {...row.getRowProps()} key={row.getRowProps().key}>
                {row.cells.map((cell) => (
                  <TableCell
                    {...cell.getCellProps()}
                    key={cell.getCellProps().key}
                  >
                    {cell.render("Cell")}
                  </TableCell>
                ))}
              </TableRow>
            );
          })}
        </TableBody>
      </MaUTable>
    </>
  );
};

export default Table;
