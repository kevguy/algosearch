import React from "react";
import CopyIcon from "../copyable/CopyIcon";
import styles from "./AddressHeader.module.scss";

const AddressHeader = ({ address }: { address: string }) => {
  return (
    <div className={styles["address-header"]}>
      <div className="sizer">
        <h3>Address Information</h3>
        <div>
          <span>{address}</span>
          <CopyIcon
            copyableText={address}
            className={styles["address-button"]}
          />
        </div>
      </div>
    </div>
  );
};

export default AddressHeader;
