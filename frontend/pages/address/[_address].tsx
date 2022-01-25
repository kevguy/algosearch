import React, { useCallback, useEffect, useState } from "react";
import axios from "axios";
import { useRouter } from "next/router";
import Layout from "../../components/layout";
import { siteName } from "../../utils/constants";
import Load from "../../components/tableloading";
import Statscard from "../../components/statscard";
import AlgoIcon from "../../components/algoicon";
import blocksTableStyles from "../../components/tableColumns/blocks.module.scss";
import statcardStyles from "../../components/statscard/Statscard.module.scss";
import {
  integerFormatter,
  microAlgosToAlgos,
  formatNumber,
} from "../../utils/stringUtils";
import Table from "../../components/table";
import { apiGetASA } from "../../utils/api";
import { Column } from "react-table";
import { IAsaMap } from "../../types/misc";
import Head from "next/head";
import { transactionsColumns } from "../../components/tableColumns/transactionsColumns";

export type DataType = {
  "amount-without-pending-rewards": number;
  "pending-rewards": number;
  rewards: number;
  status: string;
};

const Address = () => {
  const router = useRouter();
  const { _address, page } = router.query;
  const [address, setAddress] = useState("");
  const [accountTxNum, setAccountTxNum] = useState(0);
  const [accountTxns, setAccountTxns] = useState();
  const [data, setData] = useState<DataType>();
  const [loading, setLoading] = useState(true);
  const [tableLoading, setTableLoading] = useState(true);
  const [pageSize, setPageSize] = useState(15);
  const [pageCount, setPageCount] = useState(0);
  const [asaMap, setAsaMap] = useState<IAsaMap>([]);
  const [columns, setColumns] = useState<Column[]>([]);

  const getAddressData = async (address: string) => {
    await axios({
      method: "get",
      url: `${siteName}/v1/accounts/${address}?page=1&limit=10&order=desc`,
    })
      .then((response) => {
        setData(response.data);
        setLoading(false);
      })
      .catch((error) => {
        console.error(
          "Exception when querying for address information: " + error
        );
      });
  };

  const getAccountTxs = useCallback(
    async (pageIndex: number) => {
      if (!address) return;
      await axios({
        method: "get",
        url: `${siteName}/v1/transactions/acct/${address}?page=${
          pageIndex + 1
        }&limit=${pageSize}`,
      })
        .then((response) => {
          setPageCount(response.data.num_of_pages);
          setAccountTxNum(response.data.num_of_txns);
          if (response.data.items) {
            setAccountTxns(response.data.items);
          } else {
            // no account transactions with this page number, reset to first page
            router.replace({
              query: Object.assign({}, router.query, { page: "1" }),
            });
          }
        })
        .catch((error) => {
          console.error(
            "Exception when querying for address transactions: " + error
          );
        });
    },
    [address, pageSize, router]
  );

  const fetchData = useCallback(
    ({ pageIndex }) => {
      getAccountTxs(pageIndex);
    },
    [getAccountTxs]
  );

  useEffect(() => {
    if (!accountTxns) return;
    apiGetASA(accountTxns).then((result) => {
      setAsaMap(result);
    });
  }, [accountTxns]);

  useEffect(() => {
    if (!router.isReady || !_address) {
      return;
    }
    if (router.isReady && !page) {
      router.replace({
        query: Object.assign({}, router.query, { page: "1" }),
      });
    }
    setLoading(false);
    setTableLoading(false);
    setAddress(_address.toString());
    getAddressData(_address.toString());
  }, [_address, page, router]);

  useEffect(() => {
    if (address && page) {
      fetchData({ pageIndex: Number(page) - 1 });
    }
  }, [address, page, fetchData]);

  useEffect(() => {
    if (asaMap && address) {
      setColumns(transactionsColumns(asaMap, address));
    }
  }, [asaMap, address]);

  return (
    <Layout
      data={{
        address: address,
      }}
      addresspage
    >
      <Head>
        <title>AlgoSearch | Address {address.toString()}</title>
      </Head>
      <div className={`${statcardStyles["card-container"]}`}>
        <Statscard
          stat="Balance"
          value={
            <div>
              <AlgoIcon />{" "}
              {formatNumber(
                Number(
                  microAlgosToAlgos(
                    (data && data["amount-without-pending-rewards"]) || 0
                  )
                )
              )}
            </div>
          }
        />
        <Statscard
          stat="Rewards"
          value={
            loading ? (
              <Load />
            ) : (
              <div>
                <AlgoIcon />{" "}
                {data && formatNumber(Number(microAlgosToAlgos(data.rewards)))}
              </div>
            )
          }
        />
        <Statscard
          stat="Pending Rewards"
          value={
            loading ? (
              <Load />
            ) : (
              <div>
                <AlgoIcon />{" "}
                {data && microAlgosToAlgos(data["pending-rewards"])}
              </div>
            )
          }
        />
        <Statscard
          stat="Transactions"
          value={
            <div>
              {loading ? (
                <Load />
              ) : !accountTxNum ? (
                0
              ) : (
                integerFormatter.format(accountTxNum)
              )}
            </div>
          }
        />
        <Statscard
          stat="Status"
          info="Whether account is online participating in consensus on a participation node"
          value={
            loading ? (
              <Load />
            ) : (
              <div>
                {data && (
                  <>
                    <div
                      className={`status-light ${
                        data.status === "Offline"
                          ? "status-offline"
                          : "status-online"
                      }`}
                    ></div>
                    <span>{data.status}</span>
                  </>
                )}
              </div>
            )
          }
        />
      </div>
      {accountTxns && (
        <div className="table">
          <Table
            columns={columns}
            loading={tableLoading}
            data={accountTxns}
            fetchData={fetchData}
            pageCount={pageCount}
            className={`${blocksTableStyles["blocks-table"]}`}
            defaultPage={Number(page)}
          ></Table>
        </div>
      )}
    </Layout>
  );
};

export default Address;
