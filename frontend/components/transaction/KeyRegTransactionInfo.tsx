import Link from "next/link";
import { TransactionResponse } from "../../types/apiResponseTypes";
import blockStyles from "../../pages/block/Block.module.scss";
import { integerFormatter } from "../../utils/stringUtils";
import { StyledTooltip } from "../tooltip";
import { base32Encode } from "@ctrl/ts-base32";

export const KeyRegTransactionInfo = ({ tx }: { tx: TransactionResponse }) => {
  const keyRegTx = tx["keyreg-transaction"];
  if (!keyRegTx) {
    return null;
  }
  return (
    <div>
      <h4>Key Registration Information</h4>
      <div className={blockStyles["block-table"]}>
        <table cellSpacing="0">
          <tbody>
            <tr>
              <td>Mark account as participating</td>
              <td>{(!keyRegTx["non-participation"]).toString()}</td>
            </tr>
            <tr>
              <td>
                <span>Selection Participation Key</span>
                <StyledTooltip info="Public key used with the Verified Random Function (VRF) result during committee selection" />
              </td>
              <td>
                {base32Encode(
                  Buffer.from(
                    keyRegTx["selection-participation-key"],
                    "base64"
                  ),
                  undefined,
                  { padding: false }
                )}
              </td>
            </tr>
            <tr>
              <td>
                <span>Vote Participation Key</span>
                <StyledTooltip info="Participation public key used in key registration transactions" />
              </td>
              <td>
                {base32Encode(
                  Buffer.from(keyRegTx["vote-participation-key"], "base64"),
                  undefined,
                  { padding: false }
                )}
              </td>
            </tr>
            <tr>
              <td>
                <span>Vote Key Dilution</span>
                <StyledTooltip info="Number of subkeys in each batch of participation keys" />
              </td>
              <td>
                {integerFormatter.format(
                  Number(keyRegTx["vote-key-dilution"].toString())
                )}
              </td>
            </tr>
            <tr>
              <td>
                <span>Vote First Valid</span>
                <StyledTooltip info="First round this participation key is valid" />
              </td>
              <td>
                <Link href={`/block/${keyRegTx["vote-first-valid"]}`}>
                  {integerFormatter.format(
                    Number(keyRegTx["vote-first-valid"].toString())
                  )}
                </Link>
              </td>
            </tr>
            <tr>
              <td>
                <span>Vote Last Valid</span>
                <StyledTooltip info="Last round this participation key is valid" />
              </td>
              <td>
                <Link href={`/block/${keyRegTx["vote-last-valid"]}`}>
                  {integerFormatter.format(
                    Number(keyRegTx["vote-last-valid"].toString())
                  )}
                </Link>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
};
