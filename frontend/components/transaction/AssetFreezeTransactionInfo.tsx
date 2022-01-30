import Link from "next/link";
import { TransactionResponse } from "../../types/apiResponseTypes";
import blockStyles from "../../pages/block/Block.module.scss";

export const AssetFreezeTransactionInfo = ({
  tx,
}: {
  tx: TransactionResponse;
}) => {
  const afrzTx = tx["asset-freeze-transaction"];
  return (
    <div>
      <h4>Asset Freeze Information</h4>
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
              <td>Asset ID</td>
              <td>{afrzTx["asset-id"]}</td>
            </tr>
            <tr>
              <td>Address</td>
              <td>
                {afrzTx.address ? (
                  <Link href={`/address/${afrzTx.address}`}>
                    {afrzTx.address}
                  </Link>
                ) : (
                  "N/A"
                )}
              </td>
            </tr>
            <tr>
              <td>New Freeze Status</td>
              <td>{afrzTx["new-freeze-status"].toString()}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
};
