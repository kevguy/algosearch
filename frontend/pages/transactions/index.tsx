import React, { useCallback, useEffect, useState } from "react";
import axios from "axios";
import Link from "next/link";
import Layout from "../../components/layout";
import Breadcrumbs from "../../components/breadcrumbs";

import AlgoIcon from "../../components/algoicon";
import { siteName } from "../../utils/constants";
import styles from "./transactions.module.scss";
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
import Table from "../../components/table";
import Head from "next/head";
import { TransactionResponse } from "../../types/apiResponseTypes";
import { Column } from "react-table";
import Load from "../../components/tableloading";
import Statscard from "../../components/statscard";
import statscardStyles from "../../components/statscard/Statscard.module.scss";
import { apiGetASA } from "../../utils/api";
import { IAsaMap } from "../../types/misc";
import TimeAgo from "timeago-react";

const Transactions = () => {
  const [loading, setLoading] = useState(true);
  const [tableLoading, setTableLoading] = useState(true);
  const [pageSize, setPageSize] = useState(15);
  const [pageCount, setPageCount] = useState(0);
  const [page, setPage] = useState(-1);
  const [totalTransactions, setTotalTransactions] = useState(0);
  const [latestTransaction, setLatestTransaction] =
    useState<TransactionResponse>();
  const [transactions, setTransactions] = useState([]);
  const [asaMap, setAsaMap] = useState<IAsaMap>([]);

  // Update transactions based on page number
  const getTransactions = useCallback(
    (pageIndex: number) => {
      if (!latestTransaction) {
        return;
      }
      setTableLoading(true);
      axios({
        method: "get",
        url: `${siteName}/v1/transactions?latest_txn=${
          latestTransaction.id
        }&page=${pageIndex + 1}&limit=${pageSize}&order=desc`,
      })
        .then((response) => {
          console.log("txs: ", response.data);
          setPage(pageIndex);
          setPageCount(response.data.num_of_pages);
          if (pageIndex == 0) {
            setTotalTransactions(response.data.num_of_txns);
          }
          setTransactions(response.data.items);
          setTableLoading(false);
        })
        .catch((error) => {
          console.error("Exception when retrieving transactions: " + error);
        });
    },
    [latestTransaction, pageSize]
  );

  const fetchData = useCallback(
    ({ pageIndex }) => {
      console.log("latestTransaction: ", latestTransaction);
      console.log("page: ", page);
      console.log("pageIndex: ", pageIndex);
      if (latestTransaction && page != pageIndex) {
        getTransactions(pageIndex);
      }
    },
    [latestTransaction, getTransactions, page]
  );

  useEffect(() => {
    axios({
      method: "get",
      url: `${siteName}/v1/current-txn`,
    })
      .then((response) => {
        console.log("latest txn: ", response.data);
        setLoading(false);
        setLatestTransaction(response.data);
      })
      .catch((error) => {
        console.error("Error when retrieving latest statistics: " + error);
      });
  }, []);

  useEffect(() => {
    if (latestTransaction) {
      fetchData({ pageIndex: 0 });
    }
  }, [latestTransaction, fetchData]);

  useEffect(() => {
    if (!transactions) return;
    apiGetASA(transactions).then((result) => {
      console.log("results? ", result);
      setAsaMap(result);
    });
  }, [transactions]);

  // Table columns
  const columns: Column[] = [
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
      Cell: ({ value }: { value: string }) => (
        <Link href={`/address/${value}`}>{ellipseAddress(value)}</Link>
      ),
    },
    {
      Header: "To",
      accessor: "payment-transaction.receiver",
      Cell: ({
        data,
        value,
      }: {
        data: TransactionResponse[];
        value: string;
      }) => {
        const tx = data[0];
        const isAsaTransfer = tx["tx-type"] === TxType.AssetTransfer;
        const _value = isAsaTransfer
          ? tx["asset-transfer-transaction"].receiver
          : value;
        return value ? (
          <Link href={`/address/${_value}`}>{ellipseAddress(_value)}</Link>
        ) : (
          "N/A"
        );
      },
    },
    {
      Header: "Amount",
      accessor: "payment-transaction.amount",
      Cell: ({
        data,
        value,
      }: {
        data: TransactionResponse[];
        value: number;
      }) => {
        const tx = data[0];

        return (
          <span>
            {tx["tx-type"] === TxType.AssetTransfer ? (
              `${formatNumber(
                Number(
                  formatAsaAmountWithDecimal(
                    BigInt(tx["asset-transfer-transaction"].amount),
                    asaMap[tx["asset-transfer-transaction"]["asset-id"]]
                      .decimals
                  )
                )
              )} ${
                asaMap[tx["asset-transfer-transaction"]["asset-id"]].unitName
              }`
            ) : (
              <>
                <AlgoIcon /> {formatNumber(Number(microAlgosToAlgos(value)))}
              </>
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
      <div className={statscardStyles["card-container"]}>
        <Statscard
          stat="Total Transactions"
          value={
            loading ? (
              <Load />
            ) : (
              <span>{integerFormatter.format(totalTransactions)}</span>
            )
          }
        />
      </div>
      <div className="table">
        {transactions && transactions.length > 0 && (
          <Table
            columns={columns}
            data={transactions}
            fetchData={fetchData}
            pageCount={pageCount}
            loading={tableLoading}
            className={`${styles["transactions-table"]}`}
          ></Table>
        )}
      </div>
    </Layout>
  );
};

export default Transactions;
