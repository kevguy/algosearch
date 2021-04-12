import React from "react";
import { Info } from "react-feather";
import Tooltip, { TooltipProps, tooltipClasses } from "@mui/material/Tooltip";
import styles from "./Statscard.module.scss";
import Button from "@mui/material/Button";
import { styled } from "@mui/material/styles";

const StyledTooltip = styled(({ className, ...props }: TooltipProps) => (
  <Tooltip {...props} classes={{ popper: className }} />
))(({ theme }) => ({
  [`& .${tooltipClasses.arrow}`]: {
    color: "rgba(223, 240, 255, 0.7)",
  },
  [`& .${tooltipClasses.tooltip}`]: {
    backgroundColor: "rgba(223, 240, 255, 0.7)",
    color: "var(--grey-darker)",
    maxWidth: 220,
    fontFamily: "var(--font-family)",
    fontSize: "var(--font-size-xs)",
    fontWeight: "normal",
    padding: "var(--space-xs) var(--space-s)",
    border: "1px solid #dadde9",
  },
}));

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
        {info && (
          <StyledTooltip title={info} placement="top" arrow>
            <Button className={styles["tooltip-button"]}>
              <Info size="16" />
            </Button>
          </StyledTooltip>
        )}
      </div>
      {value}
    </div>
  );
};

export default Statscard;
