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
import {
  ellipseAddress,
  formatAsaAmountWithDecimal,
  formatNumber,
  getTxTypeName,
  microAlgosToAlgos,
  TxType,
} from "../../utils/stringUtils";
import { apiGetASA } from "../../utils/api";
import { TransactionResponse } from "../../types/apiResponseTypes";
import { IAsaMap } from "../../types/misc";
import Table from "../../components/table";
import Head from "next/head";
import { Row } from "react-table";

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
  const { _block } = router.query;
  const [blockNum, setBlockNum] = useState(0);
  const [data, setData] = useState<IBlockData>();
  const [transactions, setTransactions] = useState<TransactionResponse[]>();
  const [partialTxs, setPartialTxs] = useState<TransactionResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [pageSize, setPageSize] = useState(15);
  const [page, setPage] = useState(-1);
  const [pageCount, setPageCount] = useState(0);
  const [asaMap, setAsaMap] = useState<IAsaMap>([]);

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
          console.log("block: ", response.data);
          setData(response.data);
          setTransactions(response.data.transactions);
          setPageCount(Math.ceil(response.data.transactions.length / pageSize));
          setLoading(false);
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
      if (transactions && page != pageIndex) {
        setPage(pageIndex);
        const endIndex =
          transactions.length < (pageIndex + 1) * pageSize
            ? transactions.length
            : (pageIndex + 1) * pageSize;
        setPartialTxs(transactions.slice(pageIndex * pageSize, endIndex));
      }
    },
    [transactions, page, pageSize]
  );

  useEffect(() => {
    if (!_block) {
      return;
    }
    getBlock(Number(_block));
    setBlockNum(Number(_block));
  }, [_block, getBlock]);

  const columns = [
    {
      Header: "Tx ID",
      accessor: "id",
      Cell: ({ value }: { value: string }) => (
        <Link href={`/tx/${value}`}>{ellipseAddress(value)}</Link>
      ),
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
      Cell: ({ row }: { row: Row<TransactionResponse> }) => {
        const tx = row.original;
        const isAsaTransfer = tx["tx-type"] === TxType.AssetTransfer;
        const _value = isAsaTransfer
          ? tx["asset-transfer-transaction"].receiver
          : tx["payment-transaction"].receiver;
        return _value ? (
          <Link href={`/address/${_value}`}>{ellipseAddress(_value)}</Link>
        ) : (
          "N/A"
        );
      },
    },
    {
      Header: "Amount",
      accessor: "payment-transaction.amount",
      Cell: ({ row }: { row: Row<TransactionResponse> }) => {
        const tx = row.original;
        const _asaAmount =
          (tx["asset-transfer-transaction"] &&
            asaMap[tx["asset-transfer-transaction"]["asset-id"]] &&
            Number(
              formatAsaAmountWithDecimal(
                BigInt(tx["asset-transfer-transaction"].amount),
                asaMap[tx["asset-transfer-transaction"]["asset-id"]].decimals
              )
            )) ??
          0;
        const _asaUnit =
          tx["asset-transfer-transaction"] &&
          asaMap[tx["asset-transfer-transaction"]["asset-id"]] &&
          asaMap[tx["asset-transfer-transaction"]["asset-id"]].unitName;

        return (
          <span>
            {tx["tx-type"] === TxType.AssetTransfer ? (
              `${formatNumber(_asaAmount)} ${_asaUnit}`
            ) : (
              <>
                <AlgoIcon />{" "}
                {formatNumber(
                  Number(microAlgosToAlgos(tx["payment-transaction"].amount))
                )}
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
  ];

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
              <td>Round</td>
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
