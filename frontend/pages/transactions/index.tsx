import React, { useCallback, useEffect, useState } from "react";
import { useSelector } from "react-redux";
import axios from "axios";
import { Column } from "react-table";
import { useRouter } from "next/router";
import Head from "next/head";

import { siteName } from "../../utils/constants";
import tableStyles from "../../components/table/CustomTable.module.scss";
import Layout from "../../components/layout";
import Breadcrumbs from "../../components/breadcrumbs";
import Table from "../../components/table";
import { apiGetASA } from "../../utils/api";
import { IAsaMap } from "../../types/misc";
import { selectLatestTxn } from "../../features/applicationSlice";
import { transactionsColumns } from "../../components/tableColumns/transactionsColumns";
import Load from "../../components/tableloading";

const Transactions = () => {
  const router = useRouter();
  const { page } = router.query;
  const [tableLoading, setTableLoading] = useState(true);
  const [pageSize, setPageSize] = useState(15);
  const [pageCount, setPageCount] = useState(0);
  const [displayPageNum, setDisplayPageNum] = useState(0);
  const [transactions, setTransactions] = useState([]);
  const [asaMap, setAsaMap] = useState<IAsaMap>([]);
  const [columns, setColumns] = useState<Column[]>([]);
  const latestTransaction = useSelector(selectLatestTxn);
  const [lastTxForTxnsQuery, setLastTxForTxnsQuery] = useState<string>();

  // Update transactions based on page number
  const getTransactions = useCallback(
    async (pageIndex: number) => {
      if (!lastTxForTxnsQuery) return;
      await axios({
        method: "get",
        url: `${siteName}/v1/transactions?latest_txn=${lastTxForTxnsQuery}&page=${
          pageIndex + 1
        }&limit=${pageSize}&order=desc`,
      })
        .then((response) => {
          if (pageCount) {
            // if it's the first call solely for getting pageCount, keep showing loading
            setTableLoading(true);
          }
          if (response.data.items) {
            setTransactions(response.data.items);
            setTableLoading(false);
          } else {
            // no transactions with this page number, reset to first page
            setDisplayPageNum(1);
            setTableLoading(true);
          }
          setPageCount(response.data.num_of_pages);
        })
        .catch((error) => {
          setTableLoading(false);
          console.error("Exception when retrieving transactions: " + error);
        });
    },
    [lastTxForTxnsQuery, pageSize, pageCount]
  );

  const fetchData = useCallback(
    ({ pageIndex }) => {
      getTransactions(pageIndex);
    },
    [getTransactions]
  );

  useEffect(() => {
    setColumns(transactionsColumns(asaMap));
  }, [asaMap]);

  useEffect(() => {
    if (!transactions) return;
    apiGetASA(transactions).then((result) => {
      setAsaMap(result);
    });
  }, [transactions]);

  useEffect(() => {
    if (router.isReady && displayPageNum !== Number(page)) {
      if (!page) {
        router.replace({
          query: Object.assign({}, router.query, { page: 1 }),
        });
      } else {
        if (pageCount) {
          if (Number(page) > pageCount || Number(page) < 1) {
            // set URL page number param value to default 1 if URL page param is out of range
            setTableLoading(true);
            setDisplayPageNum(1);
            router.replace({
              query: Object.assign({}, router.query, { page: 1 }),
            });
          } else {
            // set URL page number param value as table page number if URL page param is in range
            setDisplayPageNum(Number(page));
          }
        }
      }
    }
    if (latestTransaction && latestTransaction !== lastTxForTxnsQuery) {
      if (
        Number(page) == 1 ||
        (page && Number(page) !== 1 && !lastTxForTxnsQuery)
      ) {
        // set last tx for query if the current page is the first page and the last tx head is not same as the latest tx
        // OR set last tx for query if the user enters the page with a page number and so the last tx head is unset
        setLastTxForTxnsQuery(latestTransaction);
      }
    }
  }, [
    latestTransaction,
    page,
    router,
    lastTxForTxnsQuery,
    pageCount,
    displayPageNum,
  ]);

  useEffect(() => {
    if (!pageCount && lastTxForTxnsQuery) {
      // first call to get pageCount
      fetchData({ pageIndex: 0 });
    }
  }, [lastTxForTxnsQuery, pageCount, fetchData]);

  return (
    <Layout>
      <Head>
        <title>AlgoSearch | Transactions</title>
      </Head>
      <Breadcrumbs
        name="Transactions"
        parentLink="/"
        parentLinkName="Home"
        currentLinkName="Transactions"
      />
      {tableLoading ? (
        <div className={tableStyles["table-loader-wrapper"]}>
          <Load />
        </div>
      ) : pageCount && displayPageNum ? (
        <div className="table">
          <Table
            columns={columns}
            loading={tableLoading}
            data={transactions}
            fetchData={fetchData}
            pageCount={pageCount}
            defaultPage={displayPageNum}
          ></Table>
        </div>
      ) : (
        <p className="center-text">No transactions</p>
      )}
    </Layout>
  );
};

export default Transactions;
