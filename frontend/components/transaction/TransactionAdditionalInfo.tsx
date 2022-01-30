import React from "react";
import Link from "next/link";
import { TransactionResponse } from "../../types/apiResponseTypes";
import {
  integerFormatter,
  microAlgosToAlgos,
  removeSpace,
} from "../../utils/stringUtils";
import styles from "../../pages/tx/TransactionDetails.module.scss";
import blockStyles from "../../pages/block/Block.module.scss";
import AlgoIcon from "../algoicon";
import algosdk from "algosdk";

const TransactionAdditionalInfo = ({
  transaction,
}: {
  transaction: TransactionResponse;
}) => (
  <div>
    <h4>Additional Information</h4>
    <div className={blockStyles["block-table"]}>
      <table cellSpacing="0">
        <thead>
          <tr>
            <th>Identifier</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>First Round</td>
            <td>
              <Link
                href={`/block/${removeSpace(
                  transaction["first-valid"].toString()
                )}`}
              >
                {integerFormatter.format(
                  Number(removeSpace(transaction["first-valid"].toString()))
                )}
              </Link>
            </td>
          </tr>
          <tr>
            <td>Last Round</td>
            <td>
              <Link
                href={`/block/${removeSpace(
                  transaction["last-valid"].toString()
                )}`}
              >
                {integerFormatter.format(
                  Number(removeSpace(transaction["last-valid"].toString()))
                )}
              </Link>
            </td>
          </tr>
          <tr>
            <td>Confirmed Round</td>
            <td>
              <Link
                href={`/block/${removeSpace(
                  transaction["confirmed-round"].toString()
                )}`}
              >
                {integerFormatter.format(
                  Number(removeSpace(transaction["confirmed-round"].toString()))
                )}
              </Link>
            </td>
          </tr>
          <tr>
            <td>Sender Rewards</td>
            <td>
              <div>
                <AlgoIcon />{" "}
                {microAlgosToAlgos(transaction["sender-rewards"] || 0)}
              </div>
            </td>
          </tr>
          <tr>
            <td>Receiver Rewards</td>
            <td>
              <div>
                <AlgoIcon />{" "}
                {microAlgosToAlgos(transaction["receiver-rewards"] || 0)}
              </div>
            </td>
          </tr>
          {transaction.signature.multisig &&
            Object.keys(transaction.signature.multisig).length > 0 && (
              <tr>
                <td className={styles["valign-top-identifier"]}>Multisig</td>
                <td className={styles["multiline-details"]}>
                  <div>Version {transaction.signature.multisig.version}</div>
                  <div>
                    Threshold: {transaction.signature.multisig.threshold}{" "}
                    signature
                    {transaction.signature.multisig.threshold! > 1 && "s"}
                  </div>
                  <h4>Subsignatures</h4>
                  {transaction.signature.multisig.subsignature?.map((sig) => {
                    const _addr = algosdk.encodeAddress(
                      Buffer.from(sig["public-key"], "base64")
                    );
                    return (
                      <Link href={`/address/${_addr}`} key={sig["public-key"]}>
                        {_addr}
                      </Link>
                    );
                  })}
                </td>
              </tr>
            )}
          <tr>
            <td>Genesis ID</td>
            <td>{transaction["genesis-id"]}</td>
          </tr>
          <tr>
            <td>Genesis Hash</td>
            <td>{transaction["genesis-hash"]}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
);

export default TransactionAdditionalInfo;
