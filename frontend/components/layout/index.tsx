import React, { useEffect, useState } from "react";
import Head from "next/head";
import * as timeago from "timeago.js";
import styles from "./Layout.module.scss";

import AddressHeader from "../addressheader";
import MainHeader from "../mainheader";
import Footer from "../footer";
import HomeHeader from "./HomeHeader";
import HomeFooter from "./HomeFooter";
import { timeAgoLocale } from "../../utils/stringUtils";

type LayoutPropsType = {
  addresspage?: boolean;
  data?: {
    address: string;
  };
  children: React.ReactNode;
  homepage?: boolean;
};

const Layout = ({ addresspage, data, homepage, children }: LayoutPropsType) => {
  const [scroll, setScroll] = useState(false);
  timeago.register("en_short", timeAgoLocale);

  useEffect(() => {
    // Scroll to top button — render behaviour
    const renderScrollTop = () => {
      let scroll_position = window.pageYOffset;
      setScroll(!scroll && scroll_position > 500);
    };
    window.addEventListener("scroll", () => renderScrollTop());
  });

  // Scroll to top button — scroll up behaviour
  const scrollToTop = () => {
    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  };

  return (
    <div className={styles.layout}>
      <Head>
        <title>AlgoSearch | Algorand Block Explorer</title>
      </Head>
      <MainHeader />
      {addresspage && data && <AddressHeader address={data.address} />}
      {homepage && <HomeHeader />}
      <div className={styles.content}>
        <div className="sizer">{children}</div>
      </div>
      {homepage && (
        <div className={styles.subfooter}>
          <HomeFooter />
        </div>
      )}
      <Footer />
      <button
        className={`${styles.scrolltop} ${scroll ? "" : styles.hiddenscroll}`}
        onClick={scrollToTop}
      >
        ➜
      </button>
    </div>
  );
};

export default Layout;
