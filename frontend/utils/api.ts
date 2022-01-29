import axios from "axios";
import { algodAddr, algodProtocol, algodToken, siteName } from "./constants";
import { IAsaResponse,  ISupply, TransactionResponse } from "../types/apiResponseTypes";
import { currencyFormatter, microAlgosToAlgos, TxType } from "./stringUtils";
import { IASAInfo, IAsaMap } from "../types/misc";
import algosdk, { LogicSigAccount } from "algosdk";

export const apiGetSupply = async () => {
  try {
    const supply = await axios({
      method: "get",
      url: `${siteName}/v1/algod/ledger/supply`,
    });
    const _onlineMoney = Number(
      microAlgosToAlgos(supply.data["online-money"])
    );
    const _results: ISupply = {
      current_round: supply.data.current_round,
      "online-money": currencyFormatter.format(_onlineMoney),
    };
    return _results;
  } catch (error) {
    console.error(
      "Error when retrieving ledger supply from Algod: " + error
    );
  }
}

export const apiGetLatestBlocks = async (currentRound: number) => {
  try {
    const latestBlocks = await axios({
      method: "get",
      url: `${siteName}/v1/rounds?latest_blk=${currentRound}&page=1&limit=10&order=desc`,
    })
    return latestBlocks.data;
  } catch(error) {
    console.log("Exception when retrieving last 10 blocks: " + error);
  }
}

export const apiGetLatestTxn = async () => {
  try {
    const latestTxn = await axios({
      method: "get",
      url: `${siteName}/v1/current-txn`,
    })
    return latestTxn.data;
  } catch(error) {
    console.log("Exception when retrieving latest transaction: " + error);
  }
}

export const apiGetASA = async (transactions: TransactionResponse[]) => {
  const dedupedAsaList = Array.from(
    new Set(
      transactions
        .filter((tx) => tx["tx-type"] === TxType.AssetTransfer)
        .map((tx) => tx["asset-transfer-transaction"]["asset-id"])
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
}

export const getLsigTEAL = async (
  lsigAc: LogicSigAccount,
  tx: TransactionResponse
) => {
  const algod = new algosdk.Algodv2(
    algodToken,
    `${algodProtocol}://${algodAddr}`,
    "4001"
  );
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
