import React, { useCallback, useEffect, useState } from "react";
import { useSelector } from "react-redux";
import axios from "axios";
import Head from "next/head";
import { useRouter } from "next/router";

import Layout from "../../components/layout";
import Breadcrumbs from "../../components/breadcrumbs";
import Statscard from "../../components/statscard";
import Load from "../../components/tableloading";
import { siteName } from "../../utils/constants";
import styles from "./blocks.module.scss";
import statscardStyles from "../../components/statscard/Statscard.module.scss";
import { integerFormatter, microAlgosToAlgos } from "../../utils/stringUtils";
import Table from "../../components/table";
import {
  selectAvgBlockTxnSpeed,
  selectWsCurrentRound,
} from "../../features/applicationSlice";
import { IBlockResponse } from "../../types/apiResponseTypes";
import { blocksColumns } from "./blocksColumns";

const Blocks = () => {
  const router = useRouter();
  const { page } = router.query;
  const [loading, setLoading] = useState(true);
  const [tableLoading, setTableLoading] = useState(true);
  const [blocks, setBlocks] = useState<IBlockResponse[]>([]);
  const [pageSize, setPageSize] = useState(15);
  const [pageCount, setPageCount] = useState(0);
  const currentRound = useSelector(selectWsCurrentRound);
  const [blockQueryRound, setBlockQueryRound] = useState<number>();
  const avgBlockTime = useSelector(selectAvgBlockTxnSpeed);

  // Get blocks based on page number
  const updateBlocks = useCallback(
    async (pageIndex: number) => {
      if (!blockQueryRound) {
        return;
      }
      await axios({
        method: "get",
        url: `${siteName}/v1/rounds?latest_blk=${blockQueryRound}&limit=${pageSize}&page=${
          pageIndex + 1
        }&order=desc`,
      })
        .then((response) => {
          setBlocks(response.data.items);
          setPageCount(response.data.num_of_pages);
        })
        .catch((error) => {
          console.log("Exception when retrieving blocks: " + error);
        });
    },
    [pageSize, blockQueryRound]
  );

  const fetchData = useCallback(
    ({ pageIndex }) => {
      updateBlocks(pageIndex);
    },
    [updateBlocks]
  );

  useEffect(() => {
    if (router.isReady && !page) {
      router.replace({
        query: Object.assign({}, router.query, { page: "1" }),
      });
    }
    if (currentRound && currentRound != blockQueryRound) {
      setLoading(false);
      setTableLoading(false);
      if (blockQueryRound != currentRound) {
        if (
          Number(page) == 1 ||
          (page && Number(page) != 1 && !blockQueryRound)
        ) {
          setBlockQueryRound(currentRound);
        }
      }
    }
  }, [currentRound, page, router, blockQueryRound]);

  return (
    <Layout>
      <Head>
        <title>AlgoSearch | Blocks</title>
      </Head>
      <Breadcrumbs
        name="Blocks"
        parentLink="/"
        parentLinkName="Home"
        currentLinkName="Blocks"
      />
      <div className={statscardStyles["card-container"]}>
        <Statscard
          stat="Latest Block"
          value={
            loading ? (
              <Load />
            ) : (
              <div>{currentRound && integerFormatter.format(currentRound)}</div>
            )
          }
        />
        <Statscard
          stat="Block Time"
          info="Average block time of last 10 blocks"
          value={loading ? <Load /> : <div>{avgBlockTime} seconds</div>}
        />
      </div>
      <div className="table">
        <div>
          {router.isReady && (
            <Table
              columns={blocksColumns}
              loading={tableLoading}
              data={blocks}
              fetchData={fetchData}
              pageCount={pageCount}
              className={styles["blocks-table"]}
              defaultPage={Number(page)}
            ></Table>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default Blocks;
