import Link from "next/link";
import AlgoIcon from "../algoicon";
import { TransactionResponse } from "../../types/apiResponseTypes";
import { IAsaMap } from "../../types/misc";
import {
  ellipseAddress,
  formatAsaAmountWithDecimal,
  formatNumber,
  microAlgosToAlgos,
  TxType,
} from "../../utils/stringUtils";

export const getInnerTxReceiver = (innerTx: TransactionResponse) => {
  const innerTxReceiver =
    innerTx["tx-type"] === TxType.AssetTransfer
      ? innerTx["asset-transfer-transaction"].receiver
      : innerTx["payment-transaction"]
      ? innerTx["payment-transaction"].receiver
      : "";
  return innerTxReceiver ? (
    <Link href={`/address/${innerTxReceiver}`}>
      {ellipseAddress(innerTxReceiver)}
    </Link>
  ) : (
    "N/A"
  );
};

export const getInnerTxCloseTo = (innerTx: TransactionResponse) => {
  const innerTxCloseTo =
    innerTx["tx-type"] === TxType.AssetTransfer
      ? innerTx["asset-transfer-transaction"]["close-to"]
      : innerTx["payment-transaction"]
      ? innerTx["payment-transaction"]["close-remainder-to"]
      : "";
  return innerTxCloseTo ? (
    <Link href={`/address/${innerTxCloseTo}`}>
      {ellipseAddress(innerTxCloseTo)}
    </Link>
  ) : (
    "N/A"
  );
};

export const getAmount = (
  txType: TxType | undefined,
  tx: TransactionResponse,
  asaMap: IAsaMap
) => {
  if (txType === TxType.AssetTransfer) {
    const axferTx = tx["asset-transfer-transaction"];
    return (
      <div>
        {axferTx &&
          asaMap[axferTx["asset-id"]] &&
          formatNumber(
            Number(
              formatAsaAmountWithDecimal(
                BigInt(axferTx.amount),
                asaMap[axferTx["asset-id"]].decimals
              ) ?? 0
            )
          )}{" "}
        {axferTx && asaMap[axferTx["asset-id"]] && (
          <Link href={`/asset/${axferTx["asset-id"]}`}>
            {asaMap[axferTx["asset-id"]].unitName}
          </Link>
        )}
      </div>
    );
  }
  if (tx["payment-transaction"]) {
    return (
      <div>
        <AlgoIcon />{" "}
        {formatNumber(
          Number(microAlgosToAlgos(tx["payment-transaction"].amount))
        )}
      </div>
    );
  }
  return <div>N/A</div>;
};

export const getCloseAmount = (
  txType: TxType | undefined,
  tx: TransactionResponse,
  asaMap: IAsaMap
) => {
  if (txType === TxType.AssetTransfer) {
    const axferTx = tx["asset-transfer-transaction"];
    return (
      <div>
        {axferTx &&
          asaMap[axferTx["asset-id"]] &&
          formatNumber(
            Number(
              formatAsaAmountWithDecimal(
                BigInt(axferTx["close-amount"]),
                asaMap[axferTx["asset-id"]].decimals
              ) ?? 0
            )
          )}{" "}
        {axferTx && asaMap[axferTx["asset-id"]] && (
          <Link href={`/asset/${axferTx["asset-id"]}`}>
            {asaMap[axferTx["asset-id"]].unitName}
          </Link>
        )}
      </div>
    );
  }
  if (tx["payment-transaction"] && tx["payment-transaction"]["close-amount"]) {
    return (
      <div>
        <AlgoIcon />{" "}
        {formatNumber(
          Number(microAlgosToAlgos(tx["payment-transaction"]["close-amount"]))
        )}
      </div>
    );
  }
  return <div>N/A</div>;
};
