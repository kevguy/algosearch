import React, { useEffect, useState } from "react";
import axios from "axios";
import moment from "moment";
import Link from "next/link";
import TimeAgo from "timeago-react";
import * as timeago from "timeago.js";
import Layout from "../components/layout";
import AlgoIcon from "../components/algoicon";
import Statscard from "../components/statscard";
import Load from "../components/tableloading";
import statscardStyles from "../components/statscard/Statscard.module.scss";
import styles from "./Home.module.scss";
import {
  currencyFormatter,
  ellipseAddress,
  integerFormatter,
  timeAgoLocale,
} from "../utils/stringUtils";
import { BigNumber } from "bignumber.js";
import Button from "@mui/material/Button";
import Table from "../components/table";
import TransactionTable from "../components/table/TransactionTable";
import { useDispatch, useSelector } from "react-redux";
import {
  getCurrentRound,
  getLatestBlocks,
  getSupply,
  selectCurrentRound,
  selectLatestBlocks,
  selectSupply,
} from "../features/applicationSlice";

const Home = () => {
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(true);
  const currentRound = useSelector(selectCurrentRound);
  const blocks = useSelector(selectLatestBlocks);
  const [price, setPrice] = useState(0);
  const [circulatingSupply, setCirculatingSupply] = useState("");
  const supply = useSelector(selectSupply);
  const dispatch = useDispatch();

  timeago.register("en_short", timeAgoLocale);

  BigNumber.config({ DECIMAL_PLACES: 2 });

  useEffect(() => {
    if (currentRound.transactions && currentRound.transactions.length > 0) {
      setTransactions([...currentRound.transactions].splice(0, 10));
    }
  }, [currentRound]);

  useEffect(() => {
    if (supply.current_round > 0) {
      dispatch(getLatestBlocks(supply.current_round));
    }
  }, [supply, dispatch]);

  const getPrice = () => {
    return axios({
      method: "get",
      url: "https://api.coingecko.com/api/v3/simple/price?ids=algorand&vs_currencies=usd",
    })
      .then((response) => {
        return response.data.algorand.usd;
      })
      .catch((error) => {
        console.error(
          "Error when retrieving Algorand price from CoinGecko: " + error
        );
      });
  };

  const getCirculatingSupply = () => {
    return axios({
      method: "get",
      url: "https://metricsapi.algorand.foundation/v1/supply/circulating?unit=algo",
    })
      .then((response) => {
        const _circulatingSupply = currencyFormatter.format(response.data);
        return _circulatingSupply;
      })
      .catch((error) => {
        console.error(
          "Error when retrieving Algorand circulating suppy: " + error
        );
      });
  };

  useEffect(() => {
    document.title = "AlgoSearch (ALGO) Blockchain Explorer";
    dispatch(getSupply());
    dispatch(getCurrentRound());
    Promise.all([getPrice(), getCirculatingSupply()]).then((results) => {
      console.log("Promise.all results: ", results);
      setPrice(results[0]);
      setCirculatingSupply(results[1] || "");
      setLoading(false);
    });
  }, [dispatch]);

  const block_columns = [
    {
      Header: "Block",
      accessor: "round",
      Cell: ({ value }: { value: number }) => {
        const _value: string = value.toString().replace(" ", "");
        return (
          <Link href={`/block/${_value}`}>
            {integerFormatter.format(Number(_value))}
          </Link>
        );
      },
    },
    {
      Header: "Proposer",
      accessor: "proposer",
      Cell: ({ value }: { value: string }) => (
        <Link href={`/address/${value}`}>{ellipseAddress(value)}</Link>
      ),
    },
    {
      Header: "# TX",
      accessor: "transactions",
      Cell: ({ value }: { value: [] }) => (
        <span className="nocolor">{value?.length || 0}</span>
      ),
    },
    {
      Header: "Time",
      accessor: "timestamp",
      Cell: ({ value }: { value: number }) => (
        <TimeAgo datetime={new Date(value * 1000)} locale="en_short"></TimeAgo>
      ),
    },
  ];
  const block_columns_id = { id: "home-latest-block-sizing" };

  return (
    <Layout homepage>
      <div className={statscardStyles["card-container"]}>
        <Statscard
          stat="Latest Round"
          value={
            loading ? (
              <Load />
            ) : (
              <Link href={`/block/${supply["current_round"]}`}>
                {integerFormatter.format(supply["current_round"])}
              </Link>
            )
          }
        />
        <Statscard
          stat="Online Stake"
          value={
            loading ? (
              <Load />
            ) : (
              <div>
                <AlgoIcon /> {supply["online-money"]}
              </div>
            )
          }
        />
        <Statscard
          stat="Circulating supply"
          value={
            loading ? (
              <Load />
            ) : (
              <div>
                <AlgoIcon /> {circulatingSupply}
              </div>
            )
          }
        />
        <Statscard
          stat="Algo Price"
          info="Powered by CoinGecko"
          value={loading ? <Load /> : <>${price}</>}
        />
      </div>
      <div className={styles["home-split"]}>
        <div className={styles["block-table"] + " addresses-table"}>
          <div>
            <span>Latest blocks</span>
            <Button>
              <Link href="/blocks">View more</Link>
            </Button>
          </div>
          <Table
            data={blocks}
            columns={block_columns}
            loading={loading}
            getProps={() => block_columns_id}
          />
        </div>
        <div className={styles["block-table"]}>
          <div>
            <span>Latest transactions</span>
            <Button>
              <Link href="/transactions">View more</Link>
            </Button>
          </div>
          <TransactionTable transactions={transactions} />
        </div>
      </div>
    </Layout>
  );
};

export default Home;
