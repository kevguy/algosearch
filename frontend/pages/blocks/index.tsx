import React, { useCallback, useEffect, useMemo, useState } from "react";
import axios from "axios";
import moment from "moment";
import Link from "next/link";
import Layout from "../../components/layout";
import Breadcrumbs from "../../components/breadcrumbs";
import Statscard from "../../components/statscard";
import AlgoIcon from "../../components/algoicon";
import Load from "../../components/tableloading";
import { siteName } from "../../utils/constants";
import styles from "./blocks.module.scss";
import statscardStyles from "../../components/statscard/Statscard.module.scss";
import {
  ellipseAddress,
  integerFormatter,
  microAlgosToAlgos,
} from "../../utils/stringUtils";
import Table from "../../components/table";
import { useDispatch, useSelector } from "react-redux";
import {
  getSupply,
  selectAvgBlockTxnSpeed,
  selectSupply,
  selectWsCurrentRound,
} from "../../features/applicationSlice";
import { IBlockResponse, IBlockRewards } from "../../types/apiResponseTypes";
import Head from "next/head";

const Blocks = () => {
  const [loading, setLoading] = useState(true);
  const [tableLoading, setTableLoading] = useState(true);
  const [blocks, setBlocks] = useState<IBlockResponse[]>([]);
  const [pageSize, setPageSize] = useState(15);
  const [page, setPage] = useState(-1);
  const [pageCount, setPageCount] = useState(0);
  const currentRound = useSelector(selectWsCurrentRound);
  const [rewardRate, setRewardRate] = useState<string | number>("");
  const avgBlockTime = useSelector(selectAvgBlockTxnSpeed);
  const dispatch = useDispatch();

  // Get blocks based on page number
  const updateBlocks = useCallback(
    async (pageIndex: number) => {
      // Use current round number to retrieve last 15 blocks
      if (!currentRound) return;
      await axios({
        method: "get",
        url: `${siteName}/v1/rounds?latest_blk=${currentRound}&limit=${pageSize}&page=${
          pageIndex + 1
        }&order=desc`,
      })
        .then((response) => {
          console.log("block rounds: ", response.data);
          setBlocks(response.data.items);
          if (pageIndex == 0) {
            const rewardRate =
              response.data.items
                .map((item: IBlockResponse) => item.rewards["rewards-rate"])
                .reduce((prev: number, curr: number) => prev + curr) /
              response.data.items.length;
            setRewardRate(microAlgosToAlgos(rewardRate));
          }
          setPage(pageIndex);
          setPageCount(response.data.num_of_pages);
        })
        .catch((error) => {
          console.log("Exception when retrieving blocks: " + error);
        });
    },
    [pageSize, currentRound]
  );

  const fetchData = useCallback(
    ({ pageIndex }) => {
      if (currentRound) {
        //&& page != pageIndex) { if user doesn't want to auto-update when on same page
        updateBlocks(pageIndex);
      }
    },
    [currentRound, updateBlocks] // page,
  );

  useEffect(() => {
    if (currentRound) {
      setLoading(false);
      setTableLoading(false);
      fetchData({ pageIndex: 0 });
    }
  }, [currentRound, fetchData]);

  useEffect(() => {
    dispatch(getSupply());
  }, [dispatch]);

  const columns = useMemo(
    () => [
      {
        Header: "Block",
        accessor: "round",
        Cell: ({ value }: { value: number }) => {
          const _value = value.toString().replace(" ", "");
          return (
            <Link href={`/block/${_value}`}>
              {integerFormatter.format(value)}
            </Link>
          );
        },
      },
      {
        Header: "Proposed by",
        accessor: "proposer",
        Cell: ({ value }: { value: string }) => (
          <Link href={`/address/${value}`}>{ellipseAddress(value)}</Link>
        ),
      },
      {
        Header: "# Tx",
        accessor: "transactions",
        Cell: ({ value }: { value: [] }) => (
          <span>{value ? integerFormatter.format(value.length) : 0}</span>
        ),
      },
      {
        Header: "Block Rewards",
        accessor: "rewards",
        Cell: ({ value }: { value: IBlockRewards }) => (
          <span>
            <AlgoIcon /> {microAlgosToAlgos(value["rewards-rate"])}
          </span>
        ),
      },
      {
        Header: "Time",
        accessor: "timestamp",
        Cell: ({ value }: { value: number }) => (
          <span>{moment.unix(value).format("D MMM YYYY, h:mm:ss")}</span>
        ),
      },
    ],
    []
  );

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
          stat="Average Block Time"
          value={loading ? <Load /> : <div>{avgBlockTime} seconds</div>}
        />
        <Statscard
          stat="Block Rewards"
          info={`Average block rewards in last ${pageSize} blocks`}
          value={
            loading || !currentRound ? (
              <Load />
            ) : (
              <div>
                <AlgoIcon /> {rewardRate}
              </div>
            )
          }
        />
      </div>
      <div className="table">
        <div>
          {blocks && blocks.length > 0 && (
            <Table
              columns={columns}
              loading={tableLoading}
              data={blocks}
              fetchData={fetchData}
              pageCount={pageCount}
              className={styles["blocks-table"]}
            ></Table>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default Blocks;
