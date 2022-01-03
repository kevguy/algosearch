import axios from "axios";
import { siteName } from "./constants";
import { ICurrentRoundResponse, ISupply } from "../types/apiResponseTypes";
import { currencyFormatter, microAlgosToAlgos } from "./stringUtils";

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

export const apiGetCurrentRound = async () => {
  try {
    const currentRound = await axios({
      method: "get",
      url: `${siteName}/v1/algod/current-round`,
    });
    const currentRoundData: ICurrentRoundResponse = currentRound.data;
    return currentRoundData;
  } catch(error) {
    console.error("Error when retrieving latest statistics: " + error);
  }
};

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
