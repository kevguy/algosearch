import React from "react";
import AlgoIcon from "../algoicon";
import CopyAddress from "./copyaddress";
import styles from "./AddressHeader.module.scss";

const AddressHeader = ({
  balance,
  address,
}: {
  balance: number;
  address: string;
}) => {
  const algoIconSize = 18;
  return (
    <div className={styles["address-header"]}>
      <div className="sizer">
        <div>
          <h3>Address Information</h3>
          <div>
            <span>{address}</span>
            <CopyAddress
              address={address}
              className={styles["address-button"]}
            />
          </div>
        </div>
        <div>
          <h4>Balance</h4>
          <div className={styles.balance}>
            <AlgoIcon width={algoIconSize} height={algoIconSize} />
            <span>{balance}</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AddressHeader;
