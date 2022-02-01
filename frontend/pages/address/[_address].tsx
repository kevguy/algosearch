import React, { useCallback, useEffect, useState } from "react";
import axios from "axios";
import { useRouter } from "next/router";
import Head from "next/head";
import Link from "next/link";
import { Column } from "react-table";

import Layout from "../../components/layout";
import { isLocal, siteName } from "../../utils/constants";
import Load from "../../components/tableloading";
import Statscard from "../../components/statscard";
import AlgoIcon from "../../components/algoicon";
import blockStyles from "../block/Block.module.scss";
import blocksTableStyles from "../../components/tableColumns/blocks.module.scss";
import statcardStyles from "../../components/statscard/Statscard.module.scss";
import {
  integerFormatter,
  microAlgosToAlgos,
  formatNumber,
} from "../../utils/stringUtils";
import Table from "../../components/table";
import { apiGetASAinAssetList } from "../../utils/api";
import { IAsaMap } from "../../types/misc";
import { transactionsColumns } from "../../components/tableColumns/transactionsColumns";
import {
  IAsaResponse,
  AccountResponse,
  TransactionResponse,
  CreatedApp,
  AccountOwnedAsset,
} from "../../types/apiResponseTypes";
import { createdAppsColumns } from "../../components/tableColumns/createdAppsColumns";
import { createdAssetsColumns } from "../../components/tableColumns/createdAssetsColumns";
import { base32Encode } from "@ctrl/ts-base32";
import { assetsColumns } from "../../components/tableColumns/assetsColumns";

const Address = () => {
  const router = useRouter();
  const { _address, page } = router.query;
  const [address, setAddress] = useState("");
  const [accountTxNum, setAccountTxNum] = useState(0);
  const [accountTxns, setAccountTxns] = useState<TransactionResponse[]>();
  const [data, setData] = useState<AccountResponse>();
  const [loading, setLoading] = useState(true);
  const [tableLoading, setTableLoading] = useState(true);
  const [pageSize, setPageSize] = useState(15);
  const [pageCount, setPageCount] = useState(0);
  const [assets, setAssets] = useState<AccountOwnedAsset[]>([]);
  const [createdAssets, setCreatedAssets] = useState<IAsaResponse[]>([]);
  const [createdApps, setCreatedApps] = useState<CreatedApp[]>([]);
  const [createdAssetsPageCount, setCreatedAssetsPageCount] = useState(0);
  const [createdAppsPageCount, setCreatedAppsPageCount] = useState(0);
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

  const fetchAssetsData = useCallback(
    ({ pageIndex }) => {
      if (data && data["assets"]) {
        const endIndex =
          data["assets"].length < (pageIndex + 1) * pageSize
            ? data["assets"].length
            : (pageIndex + 1) * pageSize;
        const newAssetsList = data["assets"].slice(
          pageIndex * pageSize,
          endIndex
        );
        setAssets(newAssetsList);
      }
    },
    [data, pageSize]
  );

  const fetchCreatedAssetsData = useCallback(
    ({ pageIndex }) => {
      if (data && data["created-assets"]) {
        const endIndex =
          data["created-assets"].length < (pageIndex + 1) * pageSize
            ? data["created-assets"].length
            : (pageIndex + 1) * pageSize;
        const newCreatedAssetsList = data["created-assets"].slice(
          pageIndex * pageSize,
          endIndex
        );
        setCreatedAssets(newCreatedAssetsList);
      }
    },
    [data, pageSize]
  );

  const fetchCreatedAppsData = useCallback(
    ({ pageIndex }) => {
      if (data && data["created-apps"]) {
        const endIndex =
          data["created-apps"].length < (pageIndex + 1) * pageSize
            ? data["created-apps"].length
            : (pageIndex + 1) * pageSize;
        const newCreatedAppsList = data["created-apps"].slice(
          pageIndex * pageSize,
          endIndex
        );
        setCreatedApps(newCreatedAppsList);
      }
    },
    [data, pageSize]
  );

  useEffect(() => {
    if (!data || !data.assets) return;
    apiGetASAinAssetList(data.assets).then((result) => {
      setAsaMap(result);
    });
  }, [data]);

  useEffect(() => {
    if (data) {
      if (data["created-assets"]) {
        setCreatedAssetsPageCount(
          Math.ceil(data["created-assets"].length / pageSize)
        );
      }
      if (data["created-apps"]) {
        setCreatedAppsPageCount(
          Math.ceil(data["created-apps"].length / pageSize)
        );
      }
    }
  }, [data, pageSize]);

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
          info="Balance including pending rewards"
          value={
            <div>
              <AlgoIcon />{" "}
              {formatNumber(
                Number(microAlgosToAlgos((data && data["amount"]) || 0))
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
        {!accountTxNum && (
          <Statscard
            stat="Transactions"
            value={loading ? <Load /> : <div>0</div>}
          />
        )}
        <Statscard
          stat="Status"
          info="Whether account is marked online to participate in consensus on a participation node"
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
        {data && data["sig-type"] && (
          <Statscard
            stat="Signature Type"
            value={
              <div>
                {data["sig-type"] === "msig"
                  ? "MultiSig"
                  : data["sig-type"] === "lsig"
                  ? "LogicSig"
                  : "Single"}
              </div>
            }
          />
        )}
        {data &&
          data.participation &&
          data.participation["selection-participation-key"] && (
            <Statscard
              stat="Selection Participation Key"
              info="Public key used with the Verified Random Function (VRF) result during committee selection"
              value={
                <div>
                  {base32Encode(
                    Buffer.from(
                      data.participation["selection-participation-key"],
                      "base64"
                    ),
                    undefined,
                    { padding: false }
                  )}
                </div>
              }
            />
          )}
        {data &&
          data.participation &&
          data.participation["vote-participation-key"] && (
            <Statscard
              stat="Vote Participation Key"
              info="Participation public key used in key registration transactions"
              value={
                <div>
                  {base32Encode(
                    Buffer.from(
                      data.participation["vote-participation-key"],
                      "base64"
                    ),
                    undefined,
                    { padding: false }
                  )}
                </div>
              }
            />
          )}
        {data &&
          data.participation &&
          !!data.participation["vote-key-dilution"] && (
            <Statscard
              stat="Vote Key Dilution"
              info="Number of subkeys in each batch of participation keys"
              value={
                <div>
                  {integerFormatter.format(
                    Number(data.participation["vote-key-dilution"].toString())
                  )}
                </div>
              }
            />
          )}
        {data &&
          data.participation &&
          !!data.participation["vote-first-valid"] && (
            <Statscard
              stat="Vote First Valid"
              info="First round this participation key is valid"
              value={
                <Link href={`/block/${data.participation["vote-first-valid"]}`}>
                  {integerFormatter.format(
                    Number(data.participation["vote-first-valid"].toString())
                  )}
                </Link>
              }
            />
          )}
        {data &&
          data.participation &&
          !!data.participation["vote-last-valid"] && (
            <Statscard
              stat="Vote Last Valid"
              info="Last round this participation key is valid"
              value={
                <Link href={`/block/${data.participation["vote-last-valid"]}`}>
                  {integerFormatter.format(
                    Number(data.participation["vote-last-valid"].toString())
                  )}
                </Link>
              }
            />
          )}
      </div>
      {isLocal && data && data.assets && (
        <div className={blockStyles["table-group"]}>
          <h3 className={blockStyles["table-header"]}>
            {data["assets"].length > 1 &&
              integerFormatter.format(
                Number(data["assets"].length.toString())
              ) + " "}
            Assets
          </h3>
          <div className={`${blockStyles["table-wrapper"]} table`}>
            <Table
              columns={assetsColumns(asaMap, address)}
              loading={tableLoading}
              data={data["assets"]}
              fetchData={fetchAssetsData}
              pageCount={createdAssetsPageCount}
              className={`${blocksTableStyles["blocks-table"]}`}
              defaultPage={1}
              changeUrlPageParamOnPageChange={false}
            ></Table>
          </div>
        </div>
      )}
      {data && data["created-assets"] && (
        <div className={blockStyles["table-group"]}>
          <h3 className={blockStyles["table-header"]}>
            {data["created-assets"].length > 1 &&
              integerFormatter.format(
                Number(data["created-assets"].length.toString())
              ) + " "}
            Created Assets
          </h3>
          <div className={`${blockStyles["table-wrapper"]} table`}>
            <Table
              columns={createdAssetsColumns(address)}
              loading={tableLoading}
              data={data["created-assets"]}
              fetchData={fetchCreatedAssetsData}
              pageCount={createdAssetsPageCount}
              className={`${blocksTableStyles["blocks-table"]}`}
              defaultPage={1}
              changeUrlPageParamOnPageChange={false}
            ></Table>
          </div>
        </div>
      )}
      {data && data["created-apps"] && (
        <div className={blockStyles["table-group"]}>
          <h3 className={blockStyles["table-header"]}>
            {data["created-apps"].length > 1 &&
              integerFormatter.format(
                Number(data["created-apps"].length.toString())
              ) + " "}
            Created Apps
          </h3>
          <div className={`${blockStyles["table-wrapper"]} table`}>
            <Table
              columns={createdAppsColumns}
              loading={tableLoading}
              data={data["created-apps"]}
              fetchData={fetchCreatedAppsData}
              pageCount={createdAppsPageCount}
              className={`${blocksTableStyles["blocks-table"]}`}
              defaultPage={1}
              changeUrlPageParamOnPageChange={false}
            ></Table>
          </div>
        </div>
      )}
      {accountTxns && (
        <div className={blockStyles["table-group"]}>
          <h3 className={blockStyles["table-header"]}>
            {accountTxns &&
              accountTxns.length > 1 &&
              integerFormatter.format(Number(accountTxns.length.toString())) +
                " "}
            Transactions
          </h3>
          <div className={`${blockStyles["table-wrapper"]} table`}>
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
        </div>
      )}
    </Layout>
  );
};

export default Address;
