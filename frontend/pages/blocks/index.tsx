import React, { useEffect, useState } from "react";
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
import { getSupply, selectSupply } from "../../features/applicationSlice";
import { IBlockResponse, IBlockRewards } from "../../types/apiResponseTypes";

const Blocks = () => {
  const [blocks, setBlocks] = useState([]);
  const [pageSize, setPageSize] = useState(25);
  const [pages, setPages] = useState(-1);
  const [loading, setLoading] = useState(true);
  const [currentRound, setCurrentRound] = useState(0);
  const [rewardRate, setRewardRate] = useState<string | number>("");
  const [avgBlockTime, setAvgBlockTime] = useState(0);
  const supply = useSelector(selectSupply);
  const dispatch = useDispatch();

  // Update page size
  const updatePageSize = (pageIndex: number, pageSize: number) => {
    setPageSize(pageSize);
    setPages(Math.ceil(currentRound / pageSize));
    updateBlocks(pageIndex);
  };

  // Update blocks based on page number
  const updateBlocks = (pageIndex: number) => {
    // Let the request headerblock be currentRound - (current page * pageSize)
    let headBlock = pageIndex * pageSize;

    // axios({
    // 	method: 'get',
    // 	url: `${siteName}/all/blocks/${headBlock + pageSize}/${pageSize}/0` // Use pageSize from state
    // }).then(response => {
    // 	setBlocks(response.data); // Set blocks to new data to render
    // }).catch(error => {
    // 	console.log("Exception when updating blocks: " + error);
    // })
  };

  useEffect(() => {
    if (supply.current_round > 0) {
      getBlocks(supply.current_round);
    }
  }, [supply, dispatch]);

  // Get initial blocks on load
  const getBlocks = (currentRound: number) => {
    // Use current round number to retrieve last 25 blocks
    setLoading(false);
    axios({
      method: "get",
      url: `${siteName}/v1/rounds?latest_blk=${currentRound}&limit=25&page=1&order=desc`,
    })
      .then((response) => {
        console.log("block rounds: ", response.data);
        setBlocks(response.data.items);
        const rewardRate =
          response.data.items
            .map((item: IBlockResponse) => item.rewards["rewards-rate"])
            .reduce((prev: number, curr: number) => prev + curr) /
          response.data.items.length;
        setRewardRate(microAlgosToAlgos(rewardRate));
        setPages(response.data.num_of_pages);
        setLoading(false);
      })
      .catch((error) => {
        console.log("Exception when retrieving last 25 blocks: " + error);
      });
  };

  useEffect(() => {
    dispatch(getSupply());
    document.title = "AlgoSearch | Blocks";
  }, []);

  const columns = React.useMemo(
    () => [
      {
        Header: "Round",
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
        Header: "Transactions",
        accessor: "transactions",
        Cell: ({ value }: { value: [] }) => (
          <span>{value ? value.length : 0}</span>
        ),
      },
      {
        Header: "Time",
        accessor: "timestamp",
        Cell: ({ value }: { value: number }) => (
          <span>{moment.unix(value).format("D MMM YYYY, h:mm:ss")}</span>
        ),
      },
      {
        Header: "Rewards Rate",
        accessor: "rewards",
        Cell: ({ value }: { value: IBlockRewards }) => (
          <span>
            <AlgoIcon /> {microAlgosToAlgos(value["rewards-rate"])}
          </span>
        ),
      },
    ],
    []
  );

  return (
    <Layout>
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
              <span>{integerFormatter.format(supply.current_round)}</span>
            )
          }
        />
        <Statscard
          stat="Average Block Time"
          value={loading ? <Load /> : <span>{avgBlockTime}s</span>}
        />
        <Statscard
          stat="Block Reward"
          info="Average block rewards in last 25 blocks"
          value={
            loading ? (
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
          <p>
            {loading
              ? "Loading"
              : integerFormatter.format(supply.current_round)}{" "}
            blocks found
          </p>
          <p>(Showing the last {pageSize} records)</p>
        </div>
        <div>
          {/* {blocks.length > 0 && <Table
						pageIndex={0}
						pages={pages}
						data={blocks}
						columns={columns}
						loading={loading}
						pageSize={pageSize}
						defaultPageSize={25}
						pageSizeOptions={[25, 50, 100]}
						onPageChange={pageIndex => updateBlocks(pageIndex)}
						onPageSizeChange={(pageSize, pageIndex) => updatePageSize(pageIndex, pageSize)}
						sortable={false}
						className="blocks-table"
						manual
					/>} */}
          {blocks.length > 0 && (
            <Table
              columns={columns}
              data={blocks}
              className={styles["blocks-table"]}
            ></Table>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default Blocks;
