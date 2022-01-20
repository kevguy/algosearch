import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { IBlockResponse, ILatestBlocksResponse, ISupply, TransactionResponse } from "../types/apiResponseTypes";
import { apiGetLatestBlocks, apiGetLatestTxn, apiGetSupply} from "../utils/api"
import {State} from "../store"

export interface IApplicationState {
  currentRound: number;
  avgBlockTxnSpeedInSec: number;
  wsCurrentRound: number;
  latestBlocks: IBlockResponse[] | undefined;
  latestTxn: string;
  supply: {
    current_round: number;
    "online-money": string;
  }
}

const initialState: IApplicationState = {
  currentRound: 0,
  avgBlockTxnSpeedInSec: 0,
  wsCurrentRound: 0,
  latestBlocks: undefined,
  latestTxn: "",
  supply: {
    current_round: 0,
    "online-money": ""
  }
};

export const getSupply = createAsyncThunk("app/getSupply", async () => {
  const response = await apiGetSupply() ?? initialState.supply
  return response;
})

export const getLatestBlocks = createAsyncThunk("app/getLatestBlocks", async (curRound: number) => {
  const response: ILatestBlocksResponse = await apiGetLatestBlocks(curRound) ?? initialState.latestBlocks
  return response;
})

export const getLatestTxn = createAsyncThunk("app/getLatestTxn", async () => {
  const response = await apiGetLatestTxn() ?? initialState.latestTxn
  return response;
})

export const applicationSlice = createSlice({
  name: 'app',
  initialState,
  reducers: {
    setAvgBlockTxnSpeed(state, action) {
      state.avgBlockTxnSpeedInSec = action.payload;
    },
    setWsCurrentRound(state, action) {
      state.wsCurrentRound = action.payload;
    },
  },
  extraReducers(builder) {
    builder
      .addCase(getLatestBlocks.fulfilled, (state, action: PayloadAction<ILatestBlocksResponse>) => {
        state.latestBlocks = action.payload && action.payload.items;
      })
      .addCase(getSupply.fulfilled, (state, action: PayloadAction<ISupply>) => {
        state.supply = action.payload;
        state.currentRound = action.payload.current_round;
      })
      .addCase(getLatestTxn.fulfilled, (state, action: PayloadAction<TransactionResponse>) => {
        state.latestTxn = action.payload.id;
      })
  }
});

export const selectCurrentRound = (state: State) => state.app.currentRound;
export const selectWsCurrentRound = (state: State) => state.app.wsCurrentRound;
export const selectLatestBlocks = (state: State) => state.app.latestBlocks;
export const selectSupply = (state: State) => state.app.supply;
export const selectAvgBlockTxnSpeed = (state: State) => state.app.avgBlockTxnSpeedInSec;
export const selectLatestTxn = (state: State) => state.app.latestTxn;

export const {
  setAvgBlockTxnSpeed,
  setWsCurrentRound,
} = applicationSlice.actions;

export default applicationSlice.reducer;
