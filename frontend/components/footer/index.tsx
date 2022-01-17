import React from "react";
import styles from "./Footer.module.scss";

const Footer = () => {
  const year = new Date().getFullYear();
  return (
    <div className={styles.footer}>
      <div>
        <p>AlgoSearch &copy; {year}</p>
      </div>
    </div>
  );
};

export default Footer;
