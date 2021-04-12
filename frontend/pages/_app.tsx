import React from "react";
import { AppProps } from "next/app";
import { Provider } from "react-redux";
import store from "../store";
import "../styles/styles.scss";

// This default export is required in a new `pages/_app.js` file.
export default function MyApp({ Component, pageProps }: AppProps) {
  return (
    <Provider store={store}>
      <Component {...pageProps} />
    </Provider>
  );
}
