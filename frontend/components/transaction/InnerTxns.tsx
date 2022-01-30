import { TransactionResponse } from "../../types/apiResponseTypes";

import styles from "../../pages/tx/TransactionDetails.module.scss";
import blockStyles from "../../pages/block/Block.module.scss";
import {
  ellipseAddress,
  getTxTypeName,
  microAlgosToAlgos,
} from "../../utils/stringUtils";
import Link from "next/link";
import {
  getAmount,
  getCloseAmount,
  getInnerTxCloseTo,
  getInnerTxReceiver,
} from "./TransactionContentComponents";
import AlgoIcon from "../algoicon";
import { IAsaMap } from "../../types/misc";

export const InnerTxns = ({
  tx,
  asaMap,
}: {
  tx: TransactionResponse;
  asaMap: IAsaMap;
}) => {
  return (
    <div>
      <h4>Inner Transactions</h4>
      <div
        className={`${blockStyles["block-table"]} ${styles["inner-txs-table"]}`}
      >
        <table cellSpacing="0">
          <thead>
            <tr>
              <th>Type</th>
              <th>Sender</th>
              <th>Receiver</th>
              <th>Amount</th>
              <th>Close To</th>
              <th>Close Amount</th>
              <th>Fee</th>
            </tr>
          </thead>
          <tbody>
            {tx["inner-txns"]!.map((innerTx, index) => (
              <tr key={index}>
                <td className={styles["normal-text"]}>
                  <h4 className="mobile-only">Type</h4>
                  {getTxTypeName(innerTx["tx-type"])}
                </td>
                <td>
                  <h4 className="mobile-only">Sender</h4>
                  <Link href={`/address/${innerTx.sender}`}>
                    {ellipseAddress(innerTx.sender)}
                  </Link>
                </td>
                <td>
                  <h4 className="mobile-only">Receiver</h4>
                  {getInnerTxReceiver(innerTx)}
                </td>
                <td>
                  <h4 className="mobile-only">Amount</h4>
                  {getAmount(innerTx["tx-type"], innerTx, asaMap)}
                </td>
                <td>
                  <h4 className="mobile-only">Close To</h4>
                  {getInnerTxCloseTo(innerTx)}
                </td>
                <td>
                  <h4 className="mobile-only">Close Amount</h4>
                  {getCloseAmount(innerTx["tx-type"], innerTx, asaMap)}
                </td>
                <td>
                  <h4 className="mobile-only">Fee</h4>
                  <AlgoIcon /> {microAlgosToAlgos(innerTx.fee)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
