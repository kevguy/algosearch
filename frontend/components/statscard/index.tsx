import React from "react";
import { Info } from "react-feather";
import ReactTooltip from "react-tooltip";
import styles from "./Statscard.module.scss";

const Statscard = ({
  stat,
  info,
  value,
}: {
  stat: string;
  info?: string;
  value: JSX.Element;
}) => {
  return (
    <div className={styles.statscard}>
      <div className={styles.title}>
        <h5>{stat}</h5>
        {info && (
          <>
            <a data-tip data-for="stats-info">
              <Info size="16" />
            </a>
            <ReactTooltip
              id="stats-info"
              type="light"
              effect="solid"
              className={styles.tooltip}
            >
              <span>{info}</span>
            </ReactTooltip>
          </>
        )}
      </div>
      {value}
    </div>
  );
};

export default Statscard;
