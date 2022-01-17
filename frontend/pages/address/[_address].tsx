import React, { useCallback, useEffect, useState } from "react";
import axios from "axios";
import moment from "moment";
import { useRouter } from "next/router";
import Link from "next/link";
import Layout from "../../components/layout";
import { siteName } from "../../utils/constants";
import Load from "../../components/tableloading";
import Statscard from "../../components/statscard";
import AlgoIcon from "../../components/algoicon";
import blocksTableStyles from "../blocks/blocks.module.scss";
import statcardStyles from "../../components/statscard/Statscard.module.scss";
import {
  getTxTypeName,
  integerFormatter,
  microAlgosToAlgos,
  formatNumber,
  removeSpace,
  TxType,
  ellipseAddress,
  formatAsaAmountWithDecimal,
} from "../../utils/stringUtils";
import TimeAgo from "timeago-react";
import Table from "../../components/table";
import { apiGetASA } from "../../utils/api";
import { TransactionResponse } from "../../types/apiResponseTypes";
import { Row } from "react-table";
import { IAsaMap } from "../../types/misc";
import Head from "next/head";

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

  const getAddressData = (address: string) => {
    axios({
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
      await axios({
        method: "get",
        url: `${siteName}/v1/transactions/acct/${address}?page=${
          pageIndex + 1
        }&limit=${pageSize}`,
      })
        .then((response) => {
          console.log("account txns data: ", response.data);
          setPageCount(response.data.num_of_pages);
          setAccountTxNum(response.data.num_of_txns);
          setAccountTxns(response.data.items);
        })
        .catch((error) => {
          console.error(
            "Exception when querying for address transactions: " + error
          );
        });
    },
    [address, pageSize]
  );
  const fetchData = useCallback(
    ({ pageIndex }) => {
      if (address) {
        getAccountTxs(pageIndex);
      }
    },
    [address, getAccountTxs]
  );

  useEffect(() => {
    if (!accountTxns) return;
    apiGetASA(accountTxns).then((result) => {
      setAsaMap(result);
    });
  }, [accountTxns]);

  useEffect(() => {
    console.log("_address: ", _address);
    console.log("page: ", page);
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

  // useEffect(() => {
  //   if (address && page) {
  //     fetchData({ pageIndex: Number(page) - 1 });
  //   }
  // }, [address, page]);

  const columns = [
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
      Header: "Tx ID",
      accessor: "id",
      Cell: ({ value }: { value: string }) => (
        <Link href={`/tx/${value}`}>{ellipseAddress(value)}</Link>
      ),
    },
    {
      Header: "From",
      accessor: "sender",
      Cell: ({ value }: { value: string }) =>
        address === value ? (
          // The address' account
          <span>{ellipseAddress(value)}</span>
        ) : (
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
          address === _value ? (
            <span>{ellipseAddress(_value)}</span>
          ) : (
            <Link href={`/address/${_value}`}>{ellipseAddress(_value)}</Link>
          )
        ) : (
          "N/A"
        );
      },
    },
    {
      Header: "Type",
      accessor: "tx-type",
      Cell: ({ value }: { value: TxType }) => (
        <span>{getTxTypeName(value)}</span>
      ),
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
    {
      Header: "Time",
      accessor: "round-time",
      Cell: ({ value }: { value: number }) => (
        <span className="nocolor">
          <TimeAgo
            datetime={new Date(moment.unix(value).toDate())}
            locale="en_short"
          />
        </span>
      ),
    },
  ];

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
          stat="Pending rewards"
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
          ></Table>
        </div>
      )}
    </Layout>
  );
};

export default Address;
