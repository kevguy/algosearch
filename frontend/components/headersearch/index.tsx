import React, { useState } from "react";
import axios from "axios";
import { siteName } from "../../utils/constants";
import styles from "./HeaderSearch.module.scss";
import { IconButton, styled } from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import { useRouter } from "next/router";

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
  const search = () => {
    const search = query ? query : "";
    axios({
      method: "get",
      url: `${siteName}/detect/${search}`,
    })
      .then((response) => {
        switch (response.data) {
          case "block":
            router.push(`/block/${search}`);
            break;
          case "transaction":
            router.push(`/tx/${search}`);
            break;
          case "address":
            router.push(`/address/${search}`);
            break;
          default:
            router.push("/error");
            break;
        }
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
      <StyledIconButton aria-label="search" size="small">
        <SearchIcon fontSize="small" />
      </StyledIconButton>
    </div>
  );
};

export default HeaderSearch;
