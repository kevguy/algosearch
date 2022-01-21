import React, { useCallback, useEffect, useState } from "react";
import axios from "axios";
import { useRouter } from "next/router";
import Link from "next/link";
import Layout from "../../components/layout";
import Breadcrumbs from "../../components/breadcrumbs";
import AlgoIcon from "../../components/algoicon";
import Load from "../../components/tableloading";
import { siteName } from "../../utils/constants";
import styles from "./Block.module.scss";
import { microAlgosToAlgos } from "../../utils/stringUtils";
import { apiGetASA } from "../../utils/api";
import { TransactionResponse } from "../../types/apiResponseTypes";
import { IAsaMap } from "../../types/misc";
import Table from "../../components/table";
import Head from "next/head";
import { Column, Row } from "react-table";
import { transactionsColumns } from "../../components/tableColumns/transactionsColumns";

interface IBlockData {
  "block-hash": string;
  "genesis-hash": string;
  "genesis-id": string;
  "previous-block-hash": string;
  proposer: string;
  rewards: {
    "fee-sink": string;
    "reward-calculation-round": number;
    "reward-level": number;
    "rewards-pool": string;
    "rewards-rate": number;
    "rewards-residue": number;
  };
  round: number;
  seed: string;
  timestamp: number;
  transactions: TransactionResponse[];
}

const Block = () => {
  const router = useRouter();
  const { _block, page } = router.query;
  const [blockNum, setBlockNum] = useState(0);
  const [data, setData] = useState<IBlockData>();
  const [transactions, setTransactions] = useState<TransactionResponse[]>();
  const [partialTxs, setPartialTxs] = useState<TransactionResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [pageSize, setPageSize] = useState(15);
  const [pageCount, setPageCount] = useState(0);
  const [asaMap, setAsaMap] = useState<IAsaMap>([]);
  const [columns, setColumns] = useState<Column[]>([]);

  useEffect(() => {
    if (!transactions) return;
    apiGetASA(transactions).then((result) => {
      setAsaMap(result);
    });
  }, [transactions]);

  const getBlock = useCallback(
    (blockNum: number) => {
      if (!blockNum) return;
      axios({
        method: "get",
        url: `${siteName}/v1/algod/rounds/${blockNum}`,
      })
        .then((response) => {
          setData(response.data);
          response.data.transactions &&
            setPageCount(
              Math.ceil(response.data.transactions.length / pageSize)
            );
          setLoading(false);
          setTransactions(response.data.transactions);
        })
        .catch((error) => {
          console.error(
            `Exception when retrieving block #${blockNum}: ${error}`
          );
        });
    },
    [pageSize]
  );

  const fetchData = useCallback(
    ({ pageIndex }) => {
      if (transactions) {
        const endIndex =
          transactions.length < (pageIndex + 1) * pageSize
            ? transactions.length
            : (pageIndex + 1) * pageSize;
        const newTxs = transactions.slice(pageIndex * pageSize, endIndex);
        setPartialTxs(newTxs);
      }
    },
    [transactions, pageSize]
  );

  useEffect(() => {
    if (!_block) {
      return;
    }
    if (router.isReady && !page) {
      router.replace({
        query: Object.assign({}, router.query, { page: "1" }),
      });
    }
    getBlock(Number(_block));
    setBlockNum(Number(_block));
  }, [_block, getBlock, page, router]);

  useEffect(() => {
    if (page && transactions) {
      if (Math.ceil(transactions.length / pageSize) < Number(page)) {
        // no account transactions with this page number, reset to first page
        const newPageNum = Math.ceil(transactions.length / pageSize);
        router.replace({
          query: Object.assign({}, router.query, {
            page: `${newPageNum}`,
          }),
        });
      } else {
        fetchData({ pageIndex: Number(page) - 1 });
      }
    }
  }, [transactions, page, pageSize, router, fetchData]);

  useEffect(() => {
    setColumns(transactionsColumns(asaMap, undefined, [1, 7]));
  }, [asaMap]);

  return (
    <Layout>
      <Head>
        <title>{`AlgoSearch | Block ${blockNum}`}</title>
      </Head>
      <Breadcrumbs
        name={`Block #${blockNum}`}
        parentLink="/blocks"
        parentLinkName="Blocks"
        currentLinkName={`Block #${blockNum}`}
      />
      <div className={styles["block-table"]}>
        <table cellSpacing="0">
          <thead>
            <tr>
              <th>Identifier</th>
              <th>Value</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Block</td>
              <td>{blockNum}</td>
            </tr>
            <tr>
              <td>Timestamp</td>
              <td>
                {loading ? (
                  <Load />
                ) : (
                  <span>
                    {data && new Date(data.timestamp * 1000).toString()}
                  </span>
                )}
              </td>
            </tr>
            <tr>
              <td>Reward Rate</td>
              <td>
                {loading ? (
                  <Load />
                ) : (
                  data && (
                    <>
                      <AlgoIcon />{" "}
                      {Number(microAlgosToAlgos(data.rewards["rewards-rate"]))}
                    </>
                  )
                )}
              </td>
            </tr>
            <tr>
              <td>Proposer</td>
              <td>
                {loading ? (
                  <Load />
                ) : (
                  data && (
                    <Link href={`/address/${data.proposer}`}>
                      {data.proposer}
                    </Link>
                  )
                )}
              </td>
            </tr>
            <tr>
              <td>Block hash</td>
              <td>{loading ? <Load /> : data && data["block-hash"]}</td>
            </tr>
            <tr>
              <td>Previous block hash</td>
              <td>
                {loading ? (
                  <Load />
                ) : (
                  <Link href={`/block/${blockNum - 1}`}>
                    {data && data["previous-block-hash"]}
                  </Link>
                )}
              </td>
            </tr>
            <tr>
              <td>Seed</td>
              <td>{loading ? <Load /> : data && data.seed}</td>
            </tr>
            <tr>
              <td>Transactions</td>
              <td>{transactions ? transactions.length : 0}</td>
            </tr>
          </tbody>
        </table>
      </div>
      {transactions && transactions.length > 0 ? (
        <div>
          <h3 className={styles["table-header"]}>
            {transactions.length > 1 && transactions.length + " "}Transactions
          </h3>
          <div className={`${styles["table-wrapper"]} table`}>
            {transactions && transactions.length > 0 && (
              <Table
                columns={columns}
                data={partialTxs}
                fetchData={fetchData}
                pageCount={pageCount}
                loading={loading}
                className={`${styles["transactions-table"]}`}
              ></Table>
            )}
          </div>
        </div>
      ) : null}
    </Layout>
  );
};

export default Block;
