import Link from "next/link";
import { TransactionResponse } from "../../types/apiResponseTypes";
import blockStyles from "../../pages/block/Block.module.scss";
import { removeSpace } from "../../utils/stringUtils";

export const AssetFreezeTransactionInfo = ({
  tx,
}: {
  tx: TransactionResponse;
}) => {
  const afrzTx = tx["asset-freeze-transaction"];
  if (!afrzTx) return null;
  return (
    <div>
      <h4>Asset Freeze Information</h4>
      <div className={blockStyles["block-table"]}>
        <table cellSpacing="0">
          <tbody>
            <tr>
              <td>Asset ID</td>
              <td>
                <Link href={`/asset/${afrzTx["asset-id"]}`}>
                  {removeSpace(afrzTx["asset-id"].toString())}
                </Link>
              </td>
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
