import React from "react";
import styles from "./Statscard.module.scss";
import { StyledTooltip } from "../tooltip";

const Statscard = ({
  stat,
  info,
  value,
}: {
  stat: string;
  info?: string | JSX.Element;
  value: JSX.Element;
}) => {
  return (
    <div className={styles.statscard}>
      <div className={styles.title}>
        <h5>{stat}</h5>
        {info && <StyledTooltip info={info} />}
      </div>
      {value}
    </div>
  );
};

export default Statscard;
