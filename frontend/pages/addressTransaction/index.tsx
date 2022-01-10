import React, { useEffect, useState } from "react";
import axios from "axios";
import Layout from "../../components/layout";
import Breadcrumbs from "../../components/breadcrumbs";
import ReactTable from "react-table-6";
import AlgoIcon from "../../components/algoicon";
import "react-table-6/react-table.css";
import { siteName } from "../../utils/constants";
import moment from "moment";
import { useRouter } from "next/router";
import Link from "next/link";

const AddressTransaction = () => {
  const router = useRouter();
  const { _address } = router.query;
  const [address, setAddress] = useState("");
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState([]);

  const getAllTransactions = (address: string) => {
    axios({
      method: "get",
      url: `${siteName}/all/addresstx/${address}`,
    }).then((response) => {
      setData(response.data);
      setLoading(false);
    });
  };

  useEffect(() => {
    if (_address) {
      document.title = `AlgoSearch | Transactions for ${_address}`;
      setAddress(_address.toString());
      getAllTransactions(_address.toString());
    }
  }, [_address]);

  // Table columns
  const columns = [
    {
      Header: "#",
      accessor: "confirmed-round",
      Cell: ({ index }: { index: number }) => (
        <span className="rownumber">{index + 1}</span>
      ),
    },
    {
      Header: "Round",
      accessor: "confirmed-round",
      Cell: ({ value }: { value: number }) => (
        <Link href={`/block/${value}`}>{value}</Link>
      ),
    },
    {
      Header: "TX ID",
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
          <span className="nocolor">{value}</span>
        ) : (
          <Link href={`/address/${value}`}>{value}</Link>
        ),
    },
    {
      Header: "",
      accessor: "sender",
      Cell: ({ value }: { value: string }) =>
        address === value ? (
          <span className="type noselect">OUT</span>
        ) : (
          <span className="type type-width-in noselect">IN</span>
        ),
    },
    {
      Header: "To",
      accessor: "payment-transaction.receiver",
      Cell: ({ value }: { value: string }) =>
        address === value ? (
          <span className="nocolor">{value}</span>
        ) : (
          <Link href={`/address/${value}`}>{value}</Link>
        ),
    },
    {
      Header: "Amount",
      accessor: "payment-transaction.amount",
      Cell: ({ value }: { value: number }) => (
        <span>
          {value / 1000000} <AlgoIcon />
        </span>
      ),
    },
    {
      Header: "Time",
      accessor: "round-time",
      Cell: ({ value }: { value: number }) => (
        <span className="nocolor">{moment.unix(value).fromNow()}</span>
      ),
    },
  ];

  return (
    <Layout>
      <Breadcrumbs
        name={`Transactions List`}
        address={address}
        parentLink={`/address/${address}`}
        parentLinkName="Address Details"
        currentLinkName={`Transactions List`}
      />
      <div className="block-table addresses-table">
        <span>
          {data.length && data.length > 0
            ? `Showing all ${data.length} transaction`
            : `Loading transactions...`}
        </span>
        <div>
          <ReactTable
            data={data}
            columns={columns}
            loading={loading}
            defaultPageSize={25}
            pageSizeOptions={[25, 50, 100]}
            sortable={false}
            className="transactions-table addresses-table-sizing"
          />
        </div>
      </div>
    </Layout>
  );
};

export default AddressTransaction;
