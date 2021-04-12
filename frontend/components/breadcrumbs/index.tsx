import Link from "next/link";
import React from "react";
import styles from "./Breadcrumbs.module.scss";

const Breadcrumbs = (props: {
  name: string;
  address?: string;
  parentLink: string;
  parentLinkName: string;
  currentLinkName: string;
}) => {
  return (
    <div
      className={`${styles.breadcrumbs} ${
        props.address && props.address !== ""
          ? styles["breadcrumbs-address-tx"]
          : null
      }`}
    >
      <div className={styles.pageTitle}>
        <h2>{props.name}</h2>
        {props.address && props.address !== "" ? (
          <span>{props.address}</span>
        ) : null}
      </div>
      <div>
        <p>
          <Link href={props.parentLink}>{props.parentLinkName}</Link>{" "}
          <span className="noselect">/</span> {props.currentLinkName}
        </p>
      </div>
    </div>
  );
};

export default Breadcrumbs;
