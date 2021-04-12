import React from "react";
import Link from "next/link";
import Layout from "../../components/layout";
import styles from "./Error.module.scss";

const FourOhFour = () => {
  return (
    <Layout>
      <div className={styles["FourOhFour"]}>
        <div className={styles["error-image"]} />
        <h1>Oops!</h1>
        <p>It looks like this page has been lost.</p>
        <p>While we try to find it, how about taking a trip back home?</p>
        <Link href="/">Return home</Link>
      </div>
    </Layout>
  );
};

export default FourOhFour;
