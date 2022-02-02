import axios from "axios";
import { algodAddr, algodProtocol, algodToken, siteName } from "./constants";
import {
  AccountOwnedAsset,
  IAsaResponse,
  ISupply,
  TransactionResponse,
} from "../types/apiResponseTypes";
import { formatNumber, microAlgosToAlgos, TxType } from "./stringUtils";
import { IASAInfo, IAsaMap } from "../types/misc";
import algosdk, { LogicSigAccount, OnApplicationComplete } from "algosdk";
import {
  Application,
  ApplicationParams,
  ApplicationStateSchema,
} from "../types/algosdkTypes";

const algod = algodToken
  ? new algosdk.Algodv2(algodToken, `${algodProtocol}://${algodAddr}`, "4001")
  : undefined;

export const apiGetSupply = async () => {
  try {
    const supply = await axios({
      method: "get",
      url: `${siteName}/v1/algod/ledger/supply`,
    });
    const _onlineMoney = Number(microAlgosToAlgos(supply.data["online-money"]));
    const _results: ISupply = {
      current_round: supply.data.current_round,
      "online-money": formatNumber(_onlineMoney),
    };
    return _results;
  } catch (error) {
    console.error("Error when retrieving ledger supply from Algod: " + error);
  }
};

export const apiGetLatestBlocks = async (currentRound: number) => {
  try {
    const latestBlocks = await axios({
      method: "get",
      url: `${siteName}/v1/rounds?latest_blk=${currentRound}&page=1&limit=10&order=desc`,
    });
    return latestBlocks.data;
  } catch (error) {
    console.log("Exception when retrieving last 10 blocks: " + error);
  }
};

export const apiGetLatestTxn = async () => {
  try {
    const latestTxn = await axios({
      method: "get",
      url: `${siteName}/v1/current-txn`,
    });
    return latestTxn.data;
  } catch (error) {
    console.log("Exception when retrieving latest transaction: " + error);
  }
};

export const apiGetASA = async (transactions: TransactionResponse[]) => {
  const dedupedAsaList = Array.from(
    new Set(
      transactions
        .filter((tx) => tx["tx-type"] === TxType.AssetTransfer)
        .map((tx) => tx["asset-transfer-transaction"]!["asset-id"])
    )
  );
  const _asaList = await Promise.all(
    dedupedAsaList.map(
      async (asaId) =>
        await axios({
          method: "get",
          url: `${siteName}/v1/algod/assets/${asaId}`,
        })
          .then((response) => {
            const result: IAsaResponse = response.data;
            const _asaInfo: IASAInfo = {
              unitName: result.params["unit-name"],
              decimals: result.params.decimals,
            };
            return _asaInfo;
          })
          .catch((error) => {
            console.error("Error when retrieving Algorand ASA");
          })
    )
  );
  const _asaMap: IAsaMap = dedupedAsaList.reduce(
    (prev, asaId, index) => ({
      ...prev,
      [asaId]: _asaList[index],
    }),
    {}
  );
  return _asaMap;
};

export const apiGetASAinAssetList = async (asas: AccountOwnedAsset[]) => {
  const asaIdList = asas.map((asa) => asa["asset-id"]);
  const _asaList = await Promise.all(
    asas.map(
      async (asa) =>
        await axios({
          method: "get",
          url: `${siteName}/v1/algod/assets/${asa["asset-id"]}`,
        })
          .then((response) => {
            const result: IAsaResponse = response.data;
            const _asaInfo: IASAInfo = {
              unitName: result.params["unit-name"],
              decimals: result.params.decimals,
            };
            return _asaInfo;
          })
          .catch((error) => {
            console.error("Error when retrieving Algorand ASA");
          })
    )
  );
  const _asaMap: IAsaMap = asaIdList.reduce(
    (prev, asaId, index) => ({
      ...prev,
      [asaId]: _asaList[index],
    }),
    {}
  );
  return _asaMap;
};

export const getLsigTEAL = async (
  lsigAc: LogicSigAccount,
  tx: TransactionResponse
) => {
  if (!algod || !tx["genesis-id"]) {
    return null;
  }
  const payTx = algosdk.makePaymentTxnWithSuggestedParamsFromObject({
    from: tx.sender,
    to: tx.sender,
    amount: tx["payment-transaction"]!.amount,
    suggestedParams: {
      fee: tx.fee,
      firstRound: tx["first-valid"],
      lastRound: tx["last-valid"],
      genesisHash: tx["genesis-hash"],
      genesisID: tx["genesis-id"],
    },
  });

  // @ts-ignore
  const dr = new algosdk.modelsv2.DryrunRequest({
    txns: [
      {
        lsig: lsigAc.get_obj_for_encoding().lsig,
        txn: payTx.get_obj_for_encoding(),
      },
    ],
  });
  const dryrunResponse = await algod.dryrun(dr).do();
  return dryrunResponse;
};

const constructAppTxForDryrun = (
  tx: TransactionResponse,
  appArgs: Uint8Array[],
  forClearState: boolean = false
) => {
  if (!tx["genesis-id"]) {
    return null;
  }
  const appInfo = tx["application-transaction"];
  if (!appInfo) return;
  const approvalProgram = Uint8Array.from(
    Buffer.from(appInfo["approval-program"]!, "base64")
  );
  const clearProgram = Uint8Array.from(
    Buffer.from(appInfo["clear-state-program"]!, "base64")
  );
  let txObj = {
    from: tx.sender,
    onComplete: forClearState
      ? OnApplicationComplete.ClearStateOC
      : OnApplicationComplete.NoOpOC,
    approvalProgram,
    clearProgram,
    appArgs,
    numLocalInts: appInfo["local-state-schema"]!["num-uint"],
    numLocalByteSlices: appInfo["local-state-schema"]!["num-byte-slice"],
    numGlobalInts: appInfo["global-state-schema"]!["num-uint"],
    numGlobalByteSlices: appInfo["global-state-schema"]!["num-byte-slice"],
    suggestedParams: {
      fee: tx.fee,
      firstRound: tx["first-valid"],
      lastRound: tx["last-valid"],
      genesisHash: tx["genesis-hash"],
      genesisID: tx["genesis-id"],
    },
  };
  if (appInfo["foreign-apps"] && appInfo["foreign-apps"].length > 0) {
    txObj = Object.assign(txObj, {
      foreignApps: appInfo["foreign-apps"],
    });
  }
  if (appInfo["foreign-assets"] && appInfo["foreign-assets"].length > 0) {
    txObj = Object.assign(txObj, {
      foreignAssets: appInfo["foreign-assets"],
    });
  }

  return algosdk.makeApplicationCreateTxnFromObject(txObj);
};

const constructAppForDryrun = (tx: TransactionResponse) => {
  const appInfo = tx["application-transaction"];
  if (!appInfo) return;
  const approvalProgram = Uint8Array.from(
    Buffer.from(appInfo["approval-program"]!, "base64")
  );
  const clearProgram = Uint8Array.from(
    Buffer.from(appInfo["clear-state-program"]!, "base64")
  );
  return new Application(
    tx["created-application-index"] || 1,
    new ApplicationParams({
      approvalProgram,
      clearStateProgram: clearProgram,
      creator: tx.sender,
      localStateSchema: new ApplicationStateSchema(
        appInfo["local-state-schema"]!["num-uint"],
        appInfo["local-state-schema"]!["num-byte-slice"]
      ),
      globalStateSchema: new ApplicationStateSchema(
        appInfo["global-state-schema"]!["num-uint"],
        appInfo["global-state-schema"]!["num-byte-slice"]
      ),
    })
  );
};

export const getAppTEAL = async (tx: TransactionResponse) => {
  const appInfo = tx["application-transaction"];
  if (
    !algod ||
    !appInfo ||
    !appInfo["approval-program"] ||
    !appInfo["clear-state-program"] ||
    !tx["created-application-index"] ||
    !appInfo["application-args"] ||
    !appInfo["local-state-schema"] ||
    !appInfo["global-state-schema"]
  ) {
    return;
  }
  const appArgs = appInfo["application-args"].map((arg) =>
    Uint8Array.from(Buffer.from(arg, "base64"))
  );
  const app = constructAppForDryrun(tx)!;
  const appTx = constructAppTxForDryrun(tx, appArgs);

  if (!appTx) {
    return null;
  }

  // @ts-ignore
  const dr = new algosdk.modelsv2.DryrunRequest({
    apps: [app],
    txns: [{ txn: appTx.get_obj_for_encoding() }],
  });
  const dryrunResponse = await algod.dryrun(dr).do();
  return dryrunResponse;
};

export const getClearStateTEAL = async (tx: TransactionResponse) => {
  const appInfo = tx["application-transaction"];
  if (
    !algod ||
    !appInfo ||
    !appInfo["approval-program"] ||
    !appInfo["clear-state-program"] ||
    !tx["created-application-index"] ||
    !appInfo["application-args"] ||
    !appInfo["local-state-schema"] ||
    !appInfo["global-state-schema"]
  ) {
    return;
  }
  const appArgs = appInfo["application-args"].map((arg) =>
    Uint8Array.from(Buffer.from(arg, "base64"))
  );
  const app = constructAppForDryrun(tx)!;
  const appTx = constructAppTxForDryrun(tx, appArgs, true);

  if (!appTx) {
    return null;
  }

  // @ts-ignore
  const dr = new algosdk.modelsv2.DryrunRequest({
    apps: [app],
    txns: [{ txn: appTx.get_obj_for_encoding() }],
  });
  const dryrunResponse = await algod.dryrun(dr).do();
  return dryrunResponse;
};
