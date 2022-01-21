import React from "react";
import Link from "next/link";
import { TransactionResponse } from "../../types/apiResponseTypes";
import {
  formatAsaAmountWithDecimal,
  formatNumber,
  integerFormatter,
  microAlgosToAlgos,
  removeSpace,
  TxType,
} from "../../utils/stringUtils";
import styles from "./TransactionDetails.module.scss";
import blockStyles from "../block/Block.module.scss";
import AlgoIcon from "../../components/algoicon";
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
          {transaction["tx-type"] === TxType.AssetConfig && (
            <>
              <tr>
                <td>Asset Name</td>
                <td>
                  {transaction["asset-config-transaction"].params.url ? (
                    <a
                      href={transaction["asset-config-transaction"].params.url}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      {transaction["asset-config-transaction"].params.name}
                    </a>
                  ) : (
                    transaction["asset-config-transaction"].params.name
                  )}
                </td>
              </tr>
              <tr>
                <td>Manager</td>
                <td>
                  {transaction["asset-config-transaction"].params.manager ? (
                    <Link
                      href={`/address/${transaction["asset-config-transaction"].params.manager}`}
                    >
                      {transaction["asset-config-transaction"].params.manager}
                    </Link>
                  ) : (
                    "N/A"
                  )}
                </td>
              </tr>
              <tr>
                <td>Reserve</td>
                <td>
                  {transaction["asset-config-transaction"].params.reserve ? (
                    <Link
                      href={`/address/${transaction["asset-config-transaction"].params.reserve}`}
                    >
                      {transaction["asset-config-transaction"].params.reserve}
                    </Link>
                  ) : (
                    "N/A"
                  )}
                </td>
              </tr>
              <tr>
                <td>Freeze</td>
                <td>
                  {transaction["asset-config-transaction"].params.freeze ? (
                    <Link
                      href={`/address/${transaction["asset-config-transaction"].params.freeze}`}
                    >
                      {transaction["asset-config-transaction"].params.freeze}
                    </Link>
                  ) : (
                    "N/A"
                  )}
                </td>
              </tr>
              <tr>
                <td>Clawback</td>
                <td>
                  {transaction["asset-config-transaction"].params.clawback ? (
                    <Link
                      href={`/address/${transaction["asset-config-transaction"].params.clawback}`}
                    >
                      {transaction["asset-config-transaction"].params.clawback}
                    </Link>
                  ) : (
                    "N/A"
                  )}
                </td>
              </tr>
              <tr>
                <td>Decimals</td>
                <td>
                  {transaction["asset-config-transaction"].params.decimals}
                </td>
              </tr>
              <tr>
                <td>Total</td>
                <td>
                  {formatNumber(
                    Number(
                      formatAsaAmountWithDecimal(
                        BigInt(
                          transaction["asset-config-transaction"].params.total
                        ),
                        transaction["asset-config-transaction"].params.decimals
                      )
                    )
                  )}{" "}
                  {transaction["asset-config-transaction"].params["unit-name"]}
                </td>
              </tr>
              <tr>
                <td>Metadata Hash</td>
                <td>
                  {
                    transaction["asset-config-transaction"].params[
                      "metadata-hash"
                    ]
                  }
                </td>
              </tr>
            </>
          )}
        </tbody>
      </table>
    </div>
  </div>
);

export default TransactionAdditionalInfo;
