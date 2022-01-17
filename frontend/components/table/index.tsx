import React, {
  PropsWithChildren,
  ReactElement,
  useCallback,
  useEffect,
  useState,
} from "react";
import { useRouter } from "next/router";
import MaUTable from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import {
  Column,
  Cell,
  useTable,
  usePagination,
  Row,
  useExpanded,
  TableOptions,
} from "react-table";
import styles from "./Table.module.scss";
import Load from "../tableloading";
import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
} from "react-feather";

export interface TableProperties<T extends Record<string, unknown>>
  extends TableOptions<T> {
  columns: Column<T>[];
  data: any;
  fetchData?: Function;
  pageCount: number;
  loading: boolean;
  className: string;
  defaultPage?: number;
}

const Table = <T extends Record<string, unknown>>(
  props: PropsWithChildren<TableProperties<T>>
): ReactElement => {
  const {
    columns,
    data,
    fetchData,
    pageCount: controlledPageCount,
    loading,
    className,
    defaultPage,
  } = props;
  const router = useRouter();
  const instance = useTable<T>(
    {
      columns,
      data,
      initialState: {
        pageIndex: defaultPage || 0,
        pageSize: 15,
      },
      manualPagination: true,
      pageCount: controlledPageCount,
    },
    useExpanded,
    usePagination
  );
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    prepareRow,
    state: { pageIndex },
    page,
    canPreviousPage,
    canNextPage,
    pageOptions,
    gotoPage,
    nextPage,
    previousPage,
  } = instance;
  const [pageIndexDisplayed, setPageIndexDisplayed] =
    useState<number>(pageIndex);

  const setPageIndex = useCallback(() => {
    if (pageIndex + 1 !== pageIndexDisplayed) {
      setPageIndexDisplayed(pageIndex + 1);
    }
  }, [pageIndex, pageIndexDisplayed]);

  const firstPageClickHandler = useCallback(() => {
    gotoPage(0);
    router.replace({
      query: { page: 1 },
    });
  }, [gotoPage, router]);

  const prevPageClickHandler = useCallback(() => {
    previousPage();
    router.replace({
      query: { page: pageIndex },
    });
  }, [pageIndex, previousPage, router]);

  const pageInputChangeHandler = useCallback(() => {
    gotoPage(pageIndexDisplayed - 1); //pageIndexDisplayed - 1);
    router.replace({
      query: { page: pageIndexDisplayed },
    });
  }, [pageIndexDisplayed, gotoPage, router]);

  const nextPageClickHandler = useCallback(() => {
    nextPage();
    router.replace({
      query: { page: pageIndexDisplayed + 1 },
    });
  }, [pageIndexDisplayed, nextPage, router]);

  const finalPageClickHandler = useCallback(() => {
    gotoPage(controlledPageCount - 1);
    router.replace({
      query: { page: controlledPageCount },
    });
  }, [controlledPageCount, gotoPage, router]);

  useEffect(() => {
    if (fetchData && pageIndex + 1 !== pageIndexDisplayed) {
      fetchData({ pageIndex });
    }
  }, [fetchData, pageIndex]);

  useEffect(() => {
    setPageIndex();
  }, [setPageIndex]);

  return (
    <>
      <MaUTable
        {...getTableProps()}
        className={`${styles["mui-table"]}${className ? " " + className : ""}`}
      >
        <TableHead>
          {headerGroups.map((headerGroup) => {
            const {
              key: headerGroupKey,
              title: headerGroupTitle,
              role: headerGroupRole,
              ...getHeaderGroupProps
            } = headerGroup.getHeaderGroupProps();
            return (
              <TableRow key={headerGroupKey ?? 0} {...getHeaderGroupProps}>
                {headerGroup.headers.map((column) => (
                  <TableCell
                    {...column.getHeaderProps()}
                    key={column.getHeaderProps().key ?? 0}
                  >
                    {column.render("Header")}
                  </TableCell>
                ))}
              </TableRow>
            );
          })}
        </TableHead>
        {!loading && (
          <TableBody
            {...getTableBodyProps()}
            className={loading ? " isLoading" : ""}
          >
            {page.map((row: Row<T>) => {
              prepareRow(row);
              return (
                <tr {...row.getRowProps()} key={row.index}>
                  {row.cells.map((cell: Cell<T>) => {
                    return (
                      <td
                        className="px-6 py-4 whitespace-no-wrap text-sm leading-5 font-medium text-gray-900"
                        key={cell.getCellProps().name}
                        {...cell.getCellProps()}
                      >
                        {cell.render("Cell")}
                      </td>
                    );
                  })}
                </tr>
              );
            })}
          </TableBody>
        )}
      </MaUTable>
      {loading && (
        <div className={styles["table-loader-wrapper"]}>
          <Load />
        </div>
      )}
      {fetchData && (
        <div className={styles["pagination"]}>
          <button onClick={firstPageClickHandler} disabled={!canPreviousPage}>
            <ChevronsLeft />
          </button>{" "}
          <button onClick={prevPageClickHandler} disabled={!canPreviousPage}>
            <ChevronLeft />
          </button>{" "}
          <span>
            Page{" "}
            <input
              type="number"
              min={1}
              value={pageIndexDisplayed}
              onChange={(e) => {
                const page = e.target.value ? Number(e.target.value) : 1;
                setPageIndexDisplayed(page);
              }}
              onBlur={pageInputChangeHandler}
              className={styles["page-input"]}
            />{" "}
            of <strong>{pageOptions.length}</strong>{" "}
          </span>
          <button onClick={nextPageClickHandler} disabled={!canNextPage}>
            <ChevronRight />
          </button>{" "}
          <button onClick={finalPageClickHandler} disabled={!canNextPage}>
            <ChevronsRight />
          </button>{" "}
        </div>
      )}
    </>
  );
};

export default Table;
