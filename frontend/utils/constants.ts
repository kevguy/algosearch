declare global {
  namespace NodeJS {
    interface ProcessEnv {
      // @ts-ignore
      NODE_ENV: "development" | "production" | "test";
      NEXT_PUBLIC_API_URL: string;
      NEXT_PUBLIC_ALGOD_PROTOCOL: string | undefined;
      NEXT_PUBLIC_ALGOD_ADDR: string | undefined;
      NEXT_PUBLIC_ALGOD_TOKEN: string | undefined;
    }
  }
}

// @ts-ignore
export const siteName = process.env.NEXT_PUBLIC_API_URL;
const wsProtocol = siteName.split("://")[0] === "https" ? "wss://" : "ws://";
export const socketEndpoint = process.env.NEXT_PUBLIC_API_URL.toString()
  .replace(/.+\/{2}/, wsProtocol)
  .concat("/ws");
export const algodProtocol = process.env.NEXT_PUBLIC_ALGOD_PROTOCOL;
export const algodAddr = process.env.NEXT_PUBLIC_ALGOD_ADDR;
export const algodToken = process.env.NEXT_PUBLIC_ALGOD_TOKEN;
export const isLocal =
  algodAddr &&
  (algodAddr.indexOf("0.0.0.0") > -1 ||
    algodAddr.indexOf("127.0.0.1") > -1 ||
    algodAddr.indexOf("localhost") > -1);
