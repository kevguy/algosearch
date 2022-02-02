import { ReactElement } from "react";
import styles from "./Copyable.module.scss";
import CopyIcon from "./CopyIcon";

const Copyable = ({
  copyableText,
  children,
}: {
  copyableText: string;
  children?: ReactElement;
}) => {
  return (
    <div
      className={`${styles["copyable-wrapper"]}${
        children ? "" : " " + styles.inline
      }`}
    >
      {children || copyableText}
      <CopyIcon copyableText={copyableText} className={styles.copy} />
    </div>
  );
};

export default Copyable;
