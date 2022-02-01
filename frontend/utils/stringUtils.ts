import BigNumber from "bignumber.js";
import * as timeago from "timeago.js";
import moment from "moment-timezone";

export function removeSpace(text: string) {
  return text.replace(" ", "");
}

export function checkBase64EqualsEmpty(text: string) {
  return Buffer.from(text, "base64").join("").replaceAll("0", "") === "";
}

export function isZeroAddress(text: string) {
  return (
    text === "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAY5HFKQ" ||
    !text
  );
}

export function ellipseAddress(address = "", width = 6): string {
  return address
    ? `${address.slice(0, width)}...${address.slice(-width)}`
    : "N/A";
}

export function microAlgosToAlgos(microAlgos: number): string | number {
  return Number.isSafeInteger(microAlgos)
    ? parseFloat((microAlgos / 1e6).toString())
    : new BigNumber(microAlgos).dividedBy(1e6).toNumber();
}

export function formatAsaAmountWithDecimal(
  asaAmount: bigint,
  decimals: number
): string {
  const singleUnit = BigInt("1" + "0".repeat(decimals));
  const wholeUnits = asaAmount / singleUnit;
  const fractionalUnits = asaAmount % singleUnit;

  return (
    wholeUnits.toString() +
    "." +
    fractionalUnits.toString().padStart(decimals, "0")
  );
}

export const timeAgoLocale: timeago.LocaleFunc = (diff, index, totalSec) => {
  // diff: the time ago / time in number;
  // index: the index of array below;
  // totalSec: total seconds between date to be formatted and today's date;
  const tz = Intl.DateTimeFormat().resolvedOptions().timeZone;
  let datetime = "";
  if (index >= 12) {
    const timestampNow = new Date().getTime();
    datetime = moment
      .tz(new Date(timestampNow - (totalSec ?? 0) * 1000), tz)
      .format("D MMM YYYY, hh:mm A z");
  }
  return [
    ["just now", "right now"],
    ["%s secs ago", "in %s secs"],
    ["1 min ago", "in 1 min"],
    ["%s mins ago", "in %s mins"],
    ["1 hour ago", "in 1 hour"],
    ["%s hours ago", "in %s hours"],
    ["1 day ago", "in 1 day"],
    ["%s days ago", "in %s days"],
    ["1 week ago", "in 1 week"],
    ["%s weeks ago", "in %s weeks"],
    ["1 month ago", "in 1 month"],
    ["%s months ago", "in %s months"],
    [`${datetime}`, `${datetime}`],
    [`${datetime}`, `${datetime}`],
  ][index] as [string, string];
};

export const currencyFormatter = new Intl.NumberFormat("en-US", {
  maximumFractionDigits: 2,
});

export const integerFormatter = new Intl.NumberFormat("en-US", {
  maximumFractionDigits: 2,
});

function roundDownSignificantDigits(number: number, decimals: number) {
  let significantDigits = parseInt(number.toExponential().split("e-")[1]) || 0;
  let decimalsUpdated = (decimals || 0) + significantDigits - 1;
  decimals = Math.min(decimalsUpdated, number.toString().length);

  return Math.floor(number * Math.pow(10, decimals)) / Math.pow(10, decimals);
}

export const formatNumber = (number: number) => {
  const numberFormatter = new Intl.NumberFormat("en-US", {
    minimumFractionDigits: 0,
    maximumFractionDigits: 10,
  });
  return numberFormatter.format(roundDownSignificantDigits(number, 10));
};

export enum TxType {
  Pay = "pay",
  KeyReg = "keyreg",
  AssetConfig = "acfg",
  AssetTransfer = "axfer",
  AssetFreeze = "afrz",
  App = "appl",
}

export enum TxTypeTxResKey {
  "payment-transaction",
  "asset-transfer-transaction",
  "asset-freeze-transaction",
  "asset-config-transaction",
  "application-transaction",
}

export const AssetTxTypes = [
  TxType.AssetConfig,
  TxType.AssetTransfer,
  TxType.AssetFreeze,
];

export const AssetTxMap = {
  [TxType.AssetTransfer]: "asset-transfer-transaction",
  [TxType.AssetFreeze]: "asset-freeze-transaction",
};

export const TxMapForAmount = {
  [TxType.Pay]: "payment-transaction",
  [TxType.KeyReg]: "payment-transaction",
  [TxType.AssetConfig]: "payment-transaction",
  [TxType.AssetTransfer]: "asset-transfer-transaction",
  [TxType.AssetFreeze]: "payment-transaction",
  [TxType.App]: "payment-transaction",
};

export const getTxTypeName = (txType: TxType) => {
  switch (txType) {
    case TxType.KeyReg:
      return "Key Reg";
    case TxType.AssetConfig:
      return "ASA Config";
    case TxType.AssetTransfer:
      return "ASA Transfer";
    case TxType.AssetFreeze:
      return "ASA Freeze";
    case TxType.App:
      return "App Call";
    case TxType.Pay:
    default:
      return "Payment";
  }
};

export function prettyPrintTEAL(sourceCode: string[]) {
  let indexOfFirstColon = -1;
  return sourceCode
    .map((line, index) => {
      let _line = line;
      if (line.indexOf(":") > -1) {
        _line = "\n" + _line;
        indexOfFirstColon = index;
      }
      if (indexOfFirstColon > -1) {
        if (_line.indexOf(":") === -1 && _line.indexOf("#") === -1) {
          _line = "    " + _line;
        }
      }
      return _line;
    })
    .join("\n")
    .trimEnd();
}
