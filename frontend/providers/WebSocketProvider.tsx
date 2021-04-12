import React, {
  createContext,
  ReactElement,
  useCallback,
  useEffect,
} from "react";
import { useDispatch } from "react-redux";
import useWebSocket, { ReadyState } from "react-use-websocket";
import {
  getLatestTxn,
  getSupply,
  setAvgBlockTxnSpeed,
  setWsCurrentRound,
} from "../features/applicationSlice";
import { IBlockResponse } from "../types/apiResponseTypes";
import { socketEndpoint } from "../utils/constants";

type ReadyStateName =
  | "connecting"
  | "open"
  | "closing"
  | "closed"
  | "uninstantiated";

type ReadyStateNameMapType = {
  [key in ReadyState]: ReadyStateName;
};

export interface IWebSocketData {
  account_ids: [] | null;
  app_ids: [] | null;
  asset_ids: [] | null;
  avg_block_txn_speed: number;
  block: IBlockResponse;
  transactions_ids: [] | null;
}

export const WebSocketContext = createContext<{
  readyState: ReadyState;
}>({ readyState: -1 });

const readyStateNameMap: ReadyStateNameMapType = {
  [ReadyState.CONNECTING]: "connecting",
  [ReadyState.OPEN]: "open",
  [ReadyState.CLOSING]: "closing",
  [ReadyState.CLOSED]: "closed",
  [ReadyState.UNINSTANTIATED]: "uninstantiated",
};

export const getConnectionStatus = (readyState: ReadyState) =>
  readyStateNameMap[readyState];

const WebSocketProvider = ({ children }: { children: ReactElement }) => {
  const { lastJsonMessage, readyState } = useWebSocket(socketEndpoint, {
    retryOnError: true,
    reconnectAttempts: 5,
    reconnectInterval: 3000,
  });
  const dispatch = useDispatch();

  const getJsonMessage = useCallback(() => {
    if (readyState == ReadyState.OPEN && lastJsonMessage !== null) {
      dispatch(setAvgBlockTxnSpeed(lastJsonMessage.avg_block_txn_speed));
      dispatch(setWsCurrentRound(lastJsonMessage.block.round));
      dispatch(getSupply());
      dispatch(getLatestTxn());
    }
  }, [lastJsonMessage, readyState, dispatch]);

  useEffect(() => {
    getJsonMessage();
  }, [getJsonMessage]);

  return (
    <WebSocketContext.Provider
      value={{
        readyState: readyState,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
};

export default WebSocketProvider;
