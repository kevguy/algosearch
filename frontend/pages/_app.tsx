import React from "react";
import { AppProps } from "next/app";
import { Provider } from "react-redux";
import store from "../store";
import WebSocketProvider from "../providers/WebSocketProvider";
import "../styles/styles.scss";

export const isBrowser = typeof window !== "undefined";

export default function MyApp({ Component, pageProps }: AppProps) {
  return (
    <Provider store={store}>
      <WebSocketProvider>
        <Component {...pageProps} />
      </WebSocketProvider>
    </Provider>
  );
}
