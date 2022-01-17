import React, { useCallback, useEffect, useState } from "react";
import { useSelector } from "react-redux";
import axios from "axios";
import { Column } from "react-table";
import { useRouter } from "next/router";
import Head from "next/head";

import { siteName } from "../../utils/constants";
import styles from "./transactions.module.scss";
import Layout from "../../components/layout";
import Breadcrumbs from "../../components/breadcrumbs";
import Table from "../../components/table";
import { apiGetASA } from "../../utils/api";
import { IAsaMap } from "../../types/misc";
import { selectLatestTxn } from "../../features/applicationSlice";
import { transactionsColumns } from "./transactionsColumns";

const Transactions = () => {
  const router = useRouter();
  const { page } = router.query;
  const [tableLoading, setTableLoading] = useState(true);
  const [pageSize, setPageSize] = useState(15);
  const [pageCount, setPageCount] = useState(0);
  const [transactions, setTransactions] = useState([]);
  const [asaMap, setAsaMap] = useState<IAsaMap>([]);
  const [columns, setColumns] = useState<Column[]>([]);
  const latestTransaction = useSelector(selectLatestTxn);
  const [lastTxForTxnsQuery, setLastTxForTxnsQuery] = useState<string>();

  // Update transactions based on page number
  const getTransactions = useCallback(
    async (pageIndex: number) => {
      if (!lastTxForTxnsQuery) {
        return;
      }
      await axios({
        method: "get",
        url: `${siteName}/v1/transactions?latest_txn=${lastTxForTxnsQuery}&page=${
          pageIndex + 1
        }&limit=${pageSize}&order=desc`,
      })
        .then((response) => {
          console.log("txs: ", response.data);
          setTableLoading(false);
          setPageCount(response.data.num_of_pages);
          setTransactions(response.data.items);
        })
        .catch((error) => {
          console.error("Exception when retrieving transactions: " + error);
        });
    },
    [lastTxForTxnsQuery, pageSize]
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
    if (router.isReady && !page) {
      router.replace({
        query: Object.assign({}, router.query, { page: "1" }),
      });
    }
    if (latestTransaction && lastTxForTxnsQuery != latestTransaction) {
      if (
        Number(page) == 1 ||
        (page && Number(page) != 1 && !lastTxForTxnsQuery)
      ) {
        // set last tx for query if the current page is the first page and the last tx head is not same as the latest tx
        // OR set last tx for query if the user enters the page with a page number and so the last tx head is unset
        setLastTxForTxnsQuery(latestTransaction);
      }
    }
  }, [latestTransaction, page, router, lastTxForTxnsQuery]);

  useEffect(() => {
    if (!transactions) return;
    apiGetASA(transactions).then((result) => {
      setAsaMap(result);
    });
  }, [transactions]);

  return (
    <Layout>
      <Head>
        <title>AlgoSearch | Transactions</title>
      </Head>

      <Breadcrumbs
        name="Transactions"
        parentLink="/"
        parentLinkName="Home"
        currentLinkName="All Transactions"
      />
      <div className="table">
        {router.isReady && (
          <Table
            columns={columns}
            data={transactions}
            fetchData={fetchData}
            pageCount={pageCount}
            loading={tableLoading}
            className={`${styles["transactions-table"]}`}
            defaultPage={Number(page)}
          ></Table>
        )}
      </div>
    </Layout>
  );
};

export default Transactions;
