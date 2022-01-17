import React, { useState } from "react";
import axios from "axios";
import { siteName } from "../../utils/constants";
import styles from "./HeaderSearch.module.scss";
import { IconButton, styled } from "@mui/material";
import { useRouter } from "next/router";
import { Search } from "react-feather";

type SearchType =
  | "acct_found"
  | "application_found"
  | "asset_found"
  | "block_hash_found"
  | "block_round_found"
  | "txn_found";

type SearchResult = {
  [key in SearchType]: boolean;
};

const StyledIconButton = styled(IconButton)(({ theme }) => ({
  width: "30px",
  height: "30px",
  fontSize: "var(--font-size-s)",
  color: "white",
  background: "var(--blue-light)",
  borderRadius: 0,
  "&:hover": {
    background: "var(--blue)",
  },
}));

const HeaderSearch = () => {
  const router = useRouter();
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(false);
  const search = () => {
    const search = query ? query : "";
    setLoading(true);
    axios({
      method: "get",
      url: `${siteName}/v1/search?key=${search}`,
    })
      .then((response) => {
        console.log("search result: ", response.data);
        const result: SearchResult = response.data;
        const typeIndex: number | null = Object.values(result)
          .map((value, index) => (value ? index : null))
          .filter((value) => value != null)[0];
        const searchType: SearchType | null = typeIndex
          ? (Object.keys(result)[typeIndex] as SearchType)
          : null;
        switch (searchType) {
          case "block_round_found":
            router.push(`/block/${search}`);
            break;
          case "txn_found":
            router.push(`/tx/${search}`);
            break;
          case "acct_found":
            router.push(`/address/${search}`);
            break;
          case "asset_found":
          //   TODO -> enable when it is implemented properly
          //   router.push(`/asset/${search}`);
          //   break;
          default:
            router.push("/error");
            break;
        }
        setLoading(false);
      })
      .catch(() => {
        router.push("/error");
      });
  };
  return (
    <div className={styles.search}>
      <input
        type="search"
        aria-label="Search by Address, Transaction ID, or Block"
        onChange={(e) => setQuery(e.target.value)}
        onKeyDown={(e) => (e.key === "Enter" ? search() : null)}
        placeholder="Search by Address / Tx ID / Block"
      />
      <StyledIconButton aria-label="search" size="small" onClick={search}>
        {loading ? (
          <div className={styles["loading-icon"]}>
            <span></span>
          </div>
        ) : (
          <Search />
        )}
      </StyledIconButton>
    </div>
  );
};

export default HeaderSearch;
