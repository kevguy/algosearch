import { Info } from "react-feather";
import Button from "@mui/material/Button";
import Tooltip, { TooltipProps, tooltipClasses } from "@mui/material/Tooltip";
import { styled } from "@mui/material/styles";
import styles from "../statscard/Statscard.module.scss";

export const CustomStyledTooltip = styled(
  ({ className, ...props }: TooltipProps) => (
    <Tooltip {...props} classes={{ popper: className }} />
  )
)(({ theme }) => ({
  [`& .${tooltipClasses.arrow}`]: {
    color: "rgba(223, 240, 255, 0.85)",
  },
  [`& .${tooltipClasses.tooltip}`]: {
    backgroundColor: "rgba(223, 240, 255, 0.85)",
    color: "var(--grey-darker)",
    maxWidth: 220,
    fontFamily: "var(--font-family)",
    fontSize: "var(--font-size-xs)",
    fontWeight: "normal",
    padding: "var(--space-xs) var(--space-s)",
    border: "1px solid #dadde9",
  },
}));

export const StyledTooltip = ({ info }: { info: string | JSX.Element }) => {
  return (
    <CustomStyledTooltip title={info} placement="top" arrow>
      <Button className={styles["tooltip-button"]}>
        <Info size="16" />
      </Button>
    </CustomStyledTooltip>
  );
};
