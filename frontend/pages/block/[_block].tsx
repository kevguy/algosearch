import React, { useEffect, useState } from "react";
import axios from "axios";
import moment from "moment";
import { useRouter } from "next/router";
import Link from "next/link";
import Layout from "../../components/layout";
import Breadcrumbs from "../../components/breadcrumbs";
import ReactTable from "react-table-6";
import "react-table-6/react-table.css";
import AlgoIcon from "../../components/algoicon";
import Load from "../../components/tableloading";
import { siteName } from "../../utils/constants";
import styles from "./Block.module.scss";
import {
  getTxTypeName,
  microAlgosToAlgos,
  TxType,
} from "../../utils/stringUtils";
import { TransactionResponse } from "../tx/[_txid]";
import { IAsaMap } from "../../components/table/TransactionTable";

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
  const [transactions, setTransactions] = useState<TransactionResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [asaMap, setAsaMap] = useState<IAsaMap>([]);

  useEffect(() => {
    if (!transactions) {
      return;
    }
    async function getAsas() {
      const dedupedAsaList = Array.from(
        new Set(
          transactions
            .filter(
              (tx: TransactionResponse) =>
                tx["tx-type"] === TxType.AssetTransfer
            )
            .map(
              (tx: TransactionResponse) =>
                tx["asset-transfer-transaction"]["asset-id"]
            )
        )
      );
      const _asaList: string[] = await Promise.all(
        dedupedAsaList.map(
          async (asaId) =>
            await axios({
              method: "get",
              url: `${siteName}/v1/algod/assets/${asaId}`,
            })
              .then((response) => {
                console.log(
                  "asa unit name?",
                  response.data.params["unit-name"]
                );
                return response.data.params["unit-name"];
              })
              .catch((error) => {
                console.error("Error when retrieving Algorand ASA");
              })
        )
      );
      const _asaMap: IAsaMap = dedupedAsaList.reduce(
        (prev: {}, asaId: number, index: number) => ({
          ...prev,
          [asaId]: _asaList[index],
        }),
        {}
      );
      if (_asaMap) {
        setAsaMap(_asaMap);
      }
      console.log("_asaMap: ", _asaMap);
    }
    getAsas();
  }, [transactions]);

  const getBlock = (blockNum: number) => {
    axios({
      method: "get",
      url: `${siteName}/v1/algod/rounds/${blockNum}`,
    })
      .then((response) => {
        console.log("block: ", response.data);
        setData(response.data);
        setTransactions(response.data.transactions);
        setLoading(false);
      })
      .catch((error) => {
        console.error(`Exception when retrieving block #${blockNum}: ${error}`);
      });
  };

  useEffect(() => {
    console.log("_block: ", _block);
    if (!_block) {
      return;
    }
    document.title = `AlgoSearch | Block ${_block}`;
    getBlock(Number(_block));
    setBlockNum(Number(_block));
  }, [_block]);

  const columns = [
    {
      Header: "TX ID",
      accessor: "id",
      Cell: ({ value }: { value: string }) => (
        <Link href={`/tx/${value}`}>{value}</Link>
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
        <Link href={`/address/${value}`}>{value}</Link>
      ),
    },
    {
      Header: "To",
      accessor: "payment-transaction.receiver",
      Cell: ({ value }: { value: string }) => (
        <Link href={`/address/${value}`}>{value}</Link>
      ),
    },
    {
      Header: "Amount",
      accessor: "payment-transaction.amount",
      Cell: ({
        original,
        value,
      }: {
        original: TransactionResponse;
        value: number;
      }) => {
        console.log("props original: ", original);
        return (
          <span>
            {original["tx-type"] === TxType.AssetTransfer ? (
              `${microAlgosToAlgos(
                original["asset-transfer-transaction"].amount
              )} ${asaMap[original["asset-transfer-transaction"]["asset-id"]]}`
            ) : (
              <>
                <AlgoIcon /> {microAlgosToAlgos(value)}
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
                  data && moment.unix(data.timestamp).format("LLLL")
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
                      {microAlgosToAlgos(data.rewards["rewards-rate"])}
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
          <div className={styles["block-table"]}>
            <ReactTable
              data={transactions}
              columns={columns}
              loading={loading}
              defaultPageSize={25}
              pageSizeOptions={[10, 25, 50]}
              sortable={false}
              className={styles["transactions-table"]}
            />
          </div>
        </div>
      ) : null}
    </Layout>
  );
};

export default Block;
