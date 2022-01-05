import React from "react";
import CopyAddress from "./copyaddress";
import styles from "./AddressHeader.module.scss";

const AddressHeader = ({ address }: { address: string }) => {
  return (
    <div className={styles["address-header"]}>
      <div className="sizer">
        <h3>Address Information</h3>
        <div>
          <span>{address}</span>
          <CopyAddress address={address} className={styles["address-button"]} />
        </div>
      </div>
    </div>
  );
};

export default AddressHeader;
