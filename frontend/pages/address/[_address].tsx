import React, { useEffect, useState } from "react";
import axios from "axios";
import moment from "moment";
import { useRouter } from "next/router";
import Link from "next/link";
import Layout from "../../components/layout";
import { siteName } from "../../utils/constants";
import Load from "../../components/tableloading";
import Statscard from "../../components/statscard";
import AlgoIcon from "../../components/algoicon";
import styles from "./Address.module.css";
import ReactTable from "react-table-6";
import txTableStyles from "../transactions/transactions.module.scss";
import "react-table-6/react-table.css";
import statcardStyles from "../../components/statscard/Statscard.module.scss";
import {
  getTxTypeName,
  integerFormatter,
  microAlgosToAlgos,
  removeSpace,
  TxType,
} from "../../utils/stringUtils";
import TimeAgo from "timeago-react";

export type DataType = {
  "amount-without-pending-rewards": number;
  "pending-rewards": number;
  rewards: number;
  status: string;
};

const Address = () => {
  const router = useRouter();
  const { _address } = router.query;
  const [address, setAddress] = useState("");
  const [accountTxNum, setAccountTxNum] = useState(0);
  const [accountTxns, setAccountTxns] = useState([]);
  const [data, setData] = useState<DataType>();
  const [txData, setTxData] = useState({});
  const [loading, setLoading] = useState(true);

  const getAddressData = (address: string) => {
    axios({
      method: "get",
      url: `${siteName}/v1/accounts/${address}?page=1&limit=10&order=desc`,
    })
      .then((response) => {
        console.log("address data: ", response.data);
        setData(response.data);
        setLoading(false);
      })
      .catch((error) => {
        console.error(
          "Exception when querying for address information: " + error
        );
      });
  };

  const getAccountTx = (address: string) => {
    axios({
      method: "get",
      url: `${siteName}/v1/transactions/acct/${address}?page=1&limit=25`,
    })
      .then((response) => {
        console.log("account txns data: ", response.data);
        setAccountTxNum(response.data.num_of_txns);
        setAccountTxns(response.data.items);
        setLoading(false);
      })
      .catch((error) => {
        console.error(
          "Exception when querying for address transactions: " + error
        );
      });
  };

  useEffect(() => {
    if (!_address) {
      return;
    }
    console.log("_address: ", _address);
    setAddress(_address.toString());
    document.title = `AlgoSearch | Address ${_address.toString()}`;
    getAddressData(_address.toString());
    getAccountTx(_address.toString());
  }, [_address]);

  const columns = [
    {
      Header: "#",
      accessor: "confirmed-round",
      Cell: ({ index }: { index: number }) => (
        <span className="rownumber">{index + 1}</span>
      ),
    },

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
        <Link href={`/tx/${value}`}>{value}</Link>
      ),
    },
    {
      Header: "From",
      accessor: "sender",
      Cell: ({ value }: { value: string }) =>
        address === value ? (
          <span>{value}</span>
        ) : (
          <Link href={`/address/${value}`}>{value}</Link>
        ),
    },
    {
      Header: "To",
      accessor: "payment-transaction.receiver",
      Cell: ({ value }: { value: string }) =>
        address === value ? (
          <span>{value}</span>
        ) : (
          <Link href={`/address/${value}`}>{value}</Link>
        ),
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
      <div className={`${statcardStyles["card-container"]}`}>
        <Statscard
          stat="Balance"
          value={
            <div>
              <AlgoIcon />{" "}
              {integerFormatter.format(
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
                {data &&
                  integerFormatter.format(
                    Number(microAlgosToAlgos(data.rewards))
                  )}
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
          stat="Status"
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
      <div className={`block-table ${styles["addresses-table"]}`}>
        <h4>
          Latest {loading || !accountTxns ? 0 : accountTxns.length} transactions{" "}
        </h4>
        <div>
          <ReactTable
            data={accountTxns}
            columns={columns}
            loading={loading}
            defaultPageSize={25}
            pageSizeOptions={[25, 50, 100]}
            sortable={false}
            className={`${txTableStyles["transactions-table"]} ${txTableStyles["addresses-table-sizing"]}`}
          />
        </div>
      </div>
    </Layout>
  );
};

export default Address;
