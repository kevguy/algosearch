import BigNumber from "bignumber.js";
import * as timeago from "timeago.js";
import moment from "moment-timezone";

export function removeSpace(text: string) {
    return text.replace(" ", "");
}
  
export function ellipseAddress(address = "", width = 6): string {
    return `${address.slice(0, width)}...${address.slice(-width)}`;
}

export function microAlgosToAlgos(microAlgos: number) : string | number {
    return Number.isSafeInteger(microAlgos) ? 
        parseFloat((microAlgos/1e6).toString())
        : new BigNumber(microAlgos).dividedBy(1e6).toNumber();
}

export const timeAgoLocale: timeago.LocaleFunc = (diff, index, totalSec) => {
    // diff: the time ago / time in number;
    // index: the index of array below;
    // totalSec: total seconds between date to be formatted and today's date;
    const tz = Intl.DateTimeFormat().resolvedOptions().timeZone;
    let datetime = "";
    if (index >= 12) {
      const timestampNow = new Date().getTime()
      datetime = moment.tz(new Date(timestampNow - (totalSec??0)*1000), tz).format("D MMM YYYY, hh:mmA z");
    }
    return [
      ['just now', 'right now'],
      ['%s secs ago', 'in %s secs'],
      ['1 min ago', 'in 1 min'],
      ['%s mins ago', 'in %s mins'],
      ['1 hour ago', 'in 1 hour'],
      ['%s hours ago', 'in %s hours'],
      ['1 day ago', 'in 1 day'],
      ['%s days ago', 'in %s days'],
      ['1 week ago', 'in 1 week'],
      ['%s weeks ago', 'in %s weeks'],
      ['1 month ago', 'in 1 month'],
      ['%s months ago', 'in %s months'],
      [`${datetime}`, `${datetime}`],
      [`${datetime}`, `${datetime}`],
    ][index] as [string, string];
};

export const currencyFormatter = new Intl.NumberFormat('en-US', {
  maximumFractionDigits: 2
});

export const integerFormatter = new Intl.NumberFormat('en-US', {
  maximumFractionDigits: 2
});

export enum TxType {
  Pay = "pay",
  KeyReg = "keyreg",
  AssetConfig = "acfg",
  AssetTransfer = "axfer",
  AssetFreeze = "afrz",
  App = "appl",
}

export const AssetTxTypes = [TxType.AssetConfig, TxType.AssetTransfer, TxType.AssetFreeze]

export const AssetTxMap = {
  [TxType.AssetTransfer]: "asset-transfer-transaction",
  [TxType.AssetFreeze]: "asset-freeze-transaction"
}

export const getTxTypeName = (txType: TxType) => {
  switch (txType) {
    case TxType.KeyReg:
      return "Key Registration";
    case TxType.AssetConfig:
      return "Asset Configuration";
    case TxType.AssetTransfer:
      return "Asset Transfer";
    case TxType.App:
      return "Application Call";
    case TxType.Pay:
    default:
      return "Payment";
  }
};
