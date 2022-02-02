import React, { useCallback, useEffect, useState } from "react";
import axios from "axios";
import { useRouter } from "next/router";
import Link from "next/link";
import Layout from "../../components/layout";
import Breadcrumbs from "../../components/breadcrumbs";
import { siteName } from "../../utils/constants";
import styles from "../block/Block.module.scss";
import {
  formatAsaAmountWithDecimal,
  formatNumber,
} from "../../utils/stringUtils";
import { IAsaResponse } from "../../types/apiResponseTypes";
import Head from "next/head";

const Asset = () => {
  const router = useRouter();
  const { _asset } = router.query;
  const [assetId, setAssetId] = useState(0);
  const [data, setData] = useState<IAsaResponse>();
  const [loading, setLoading] = useState(true);

  const getASA = useCallback(() => {
    if (assetId != 0) {
      axios({
        method: "get",
        url: `${siteName}/v1/algod/assets/${assetId}`,
      })
        .then((response) => {
          setData(response.data);
          setLoading(false);
        })
        .catch((error) => {
          console.error(`Exception when retrieving ASA #${assetId}: ${error}`);
        });
    }
  }, [assetId]);

  useEffect(() => {
    getASA();
  }, [getASA]);

  useEffect(() => {
    if (!_asset) {
      return;
    }
    setAssetId(Number(_asset));
  }, [_asset]);

  return (
    <Layout>
      <Head>
        <title>{`AlgoSearch | Asset ${assetId}`}</title>
      </Head>
      <Breadcrumbs
        name={`Asset #${assetId}`}
        parentLink="/"
        parentLinkName="Home"
        currentLinkName={`Asset #${assetId}`}
      />
      <div className={styles["block-table"]}>
        <table cellSpacing="0">
          <tbody>
            <tr>
              <td>Asset ID</td>
              <td>{assetId}</td>
            </tr>
            <tr>
              <td>Asset Name</td>
              <td>{data?.params.name}</td>
            </tr>
            <tr>
              <td>URL</td>
              <td>
                {data?.params.url ? (
                  <a
                    href={
                      data.params.url.indexOf("https://") === -1
                        ? `https:/${data.params.url}`
                        : data.params.url
                    }
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    {data.params.url}
                  </a>
                ) : (
                  "N/A"
                )}
              </td>
            </tr>
            <tr>
              <td>Decimals</td>
              <td>{data?.params.decimals}</td>
            </tr>
            {data && (
              <tr>
                <td>Creator</td>
                <td>
                  <Link href={`/address/${data.params.creator}`}>
                    {data.params.creator}
                  </Link>
                </td>
              </tr>
            )}
            {data?.params.manager && (
              <tr>
                <td>Manager</td>
                <td>
                  <Link href={`/address/${data.params.manager}`}>
                    {data.params.manager}
                  </Link>
                </td>
              </tr>
            )}
            <tr>
              <td>Reserve Account</td>
              <td>
                {data?.params.reserve ? (
                  <Link href={`/address/${data.params.reserve}`}>
                    {data.params.reserve}
                  </Link>
                ) : (
                  "N/A"
                )}
              </td>
            </tr>
            <tr>
              <td>Freeze Account</td>
              <td>
                {data?.params.freeze ? (
                  <Link href={`/address/${data.params.freeze}`}>
                    {data?.params.freeze}
                  </Link>
                ) : (
                  "N/A"
                )}
              </td>
            </tr>
            {data && (
              <tr>
                <td>Total Supply</td>
                <td>
                  {formatNumber(
                    Number(
                      formatAsaAmountWithDecimal(
                        BigInt(data.params.total),
                        data.params.decimals
                      )
                    )
                  )}{" "}
                  {data.params["unit-name"]}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </Layout>
  );
};

export default Asset;
