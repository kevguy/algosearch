import Link from "next/link";
import { TransactionResponse } from "../../types/apiResponseTypes";
import {
  formatAsaAmountWithDecimal,
  formatNumber,
} from "../../utils/stringUtils";
import blockStyles from "../../pages/block/Block.module.scss";

export const AssetConfigTransactionInfo = ({
  tx,
}: {
  tx: TransactionResponse;
}) => {
  if (!tx["asset-config-transaction"]) return null;
  return (
    <div>
      <h4>Asset Configuration Information</h4>
      <div className={blockStyles["block-table"]}>
        <table cellSpacing="0">
          <tbody>
            <tr>
              <td>Asset Name</td>
              <td>
                {tx["asset-config-transaction"].params.url ? (
                  <a
                    href={tx["asset-config-transaction"].params.url}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    {tx["asset-config-transaction"].params.name}
                  </a>
                ) : (
                  tx["asset-config-transaction"].params.name
                )}
              </td>
            </tr>
            <tr>
              <td>Creator</td>
              <td>
                {tx["asset-config-transaction"].params.creator ? (
                  <Link
                    href={`/address/${tx["asset-config-transaction"].params.creator}`}
                  >
                    {tx["asset-config-transaction"].params.creator}
                  </Link>
                ) : (
                  "N/A"
                )}
              </td>
            </tr>
            <tr>
              <td>Manager</td>
              <td>
                {tx["asset-config-transaction"].params.manager ? (
                  <Link
                    href={`/address/${tx["asset-config-transaction"].params.manager}`}
                  >
                    {tx["asset-config-transaction"].params.manager}
                  </Link>
                ) : (
                  "N/A"
                )}
              </td>
            </tr>
            <tr>
              <td>Reserve</td>
              <td>
                {tx["asset-config-transaction"].params.reserve ? (
                  <Link
                    href={`/address/${tx["asset-config-transaction"].params.reserve}`}
                  >
                    {tx["asset-config-transaction"].params.reserve}
                  </Link>
                ) : (
                  "N/A"
                )}
              </td>
            </tr>
            <tr>
              <td>Freeze</td>
              <td>
                {tx["asset-config-transaction"].params.freeze ? (
                  <Link
                    href={`/address/${tx["asset-config-transaction"].params.freeze}`}
                  >
                    {tx["asset-config-transaction"].params.freeze}
                  </Link>
                ) : (
                  "N/A"
                )}
              </td>
            </tr>
            <tr>
              <td>Clawback</td>
              <td>
                {tx["asset-config-transaction"].params.clawback ? (
                  <Link
                    href={`/address/${tx["asset-config-transaction"].params.clawback}`}
                  >
                    {tx["asset-config-transaction"].params.clawback}
                  </Link>
                ) : (
                  "N/A"
                )}
              </td>
            </tr>
            <tr>
              <td>Decimals</td>
              <td>{tx["asset-config-transaction"].params.decimals}</td>
            </tr>
            <tr>
              <td>Total</td>
              <td>
                {formatNumber(
                  Number(
                    formatAsaAmountWithDecimal(
                      BigInt(tx["asset-config-transaction"].params.total),
                      tx["asset-config-transaction"].params.decimals
                    )
                  )
                )}{" "}
                {tx["asset-config-transaction"].params["unit-name"]}
              </td>
            </tr>
            <tr>
              <td>Metadata Hash</td>
              <td>{tx["asset-config-transaction"].params["metadata-hash"]}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
};
