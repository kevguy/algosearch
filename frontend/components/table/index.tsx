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
import styles from "./CustomTable.module.scss";
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
  className?: string;
  defaultPage?: number;
  changeUrlPageParamOnPageChange?: boolean;
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
    changeUrlPageParamOnPageChange = true,
  } = props;
  const router = useRouter();
  const instance = useTable<T>(
    {
      columns,
      data,
      initialState: {
        pageIndex: defaultPage ? defaultPage - 1 : 0,
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
  const [pageIndexDisplayed, setPageIndexDisplayed] = useState<number>(
    pageIndex + 1
  );
  const firstPageClickHandler = useCallback(() => {
    gotoPage(0);
    if (changeUrlPageParamOnPageChange) {
      router.replace({
        query: Object.assign({}, router.query, { page: 1 }),
      });
    }
    setPageIndexDisplayed(1);
  }, [gotoPage, router, changeUrlPageParamOnPageChange]);

  const prevPageClickHandler = useCallback(() => {
    previousPage();
    if (changeUrlPageParamOnPageChange) {
      router.replace({
        query: Object.assign({}, router.query, {
          page: pageIndexDisplayed - 1,
        }),
      });
    }
    setPageIndexDisplayed(pageIndexDisplayed - 1);
  }, [
    previousPage,
    pageIndexDisplayed,
    router,
    changeUrlPageParamOnPageChange,
  ]);

  const pageInputChangeHandler = useCallback(() => {
    if (pageIndexDisplayed) {
      if (pageIndexDisplayed <= pageOptions.length) {
        if (changeUrlPageParamOnPageChange) {
          router.replace({
            query: Object.assign({}, router.query, {
              page: pageIndexDisplayed,
            }),
          });
        }
        gotoPage(pageIndexDisplayed - 1);
      } else {
        if (changeUrlPageParamOnPageChange) {
          router.replace({
            query: Object.assign({}, router.query, {
              page: pageOptions.length,
            }),
          });
        }
        setPageIndexDisplayed(pageOptions.length);
        gotoPage(pageOptions.length - 1);
      }
    }
  }, [
    pageIndexDisplayed,
    gotoPage,
    pageOptions,
    router,
    changeUrlPageParamOnPageChange,
  ]);

  const nextPageClickHandler = useCallback(() => {
    nextPage();
    if (changeUrlPageParamOnPageChange) {
      router.replace({
        query: Object.assign({}, router.query, {
          page: pageIndexDisplayed + 1,
        }),
      });
    }
    setPageIndexDisplayed(pageIndexDisplayed + 1);
  }, [pageIndexDisplayed, nextPage, router, changeUrlPageParamOnPageChange]);

  const finalPageClickHandler = useCallback(() => {
    gotoPage(controlledPageCount - 1);
    if (changeUrlPageParamOnPageChange) {
      router.replace({
        query: Object.assign({}, router.query, { page: controlledPageCount }),
      });
    }
    setPageIndexDisplayed(controlledPageCount);
  }, [controlledPageCount, gotoPage, router, changeUrlPageParamOnPageChange]);

  useEffect(() => {
    if (!fetchData) {
      return;
    }
    if (changeUrlPageParamOnPageChange) {
      if (
        defaultPage === pageIndexDisplayed &&
        pageIndex + 1 === pageIndexDisplayed
      ) {
        // only fetch when page is set correct across the variables
        fetchData({
          pageIndex,
        });
      }
    } else {
      if (pageIndexDisplayed) {
        fetchData({
          pageIndex: pageIndexDisplayed - 1,
        });
      }
    }
  }, [
    fetchData,
    pageIndex,
    defaultPage,
    router,
    pageIndexDisplayed,
    changeUrlPageParamOnPageChange,
  ]);

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
      {fetchData && !loading && (
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
                const page =
                  e.target.value && Number(e.target.value) > 0
                    ? Number(e.target.value)
                    : 1;
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
