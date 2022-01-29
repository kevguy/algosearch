declare global {
  namespace NodeJS {
    interface ProcessEnv {
      // @ts-ignore
      NODE_ENV: "development" | "production" | "test";
      NEXT_PUBLIC_API_URL: string;
      NEXT_PUBLIC_ALGOD_PROTOCOL: string;
      NEXT_PUBLIC_ALGOD_ADDR: string;
      NEXT_PUBLIC_ALGOD_TOKEN: string;
    }
  }
}

// @ts-ignore
const wsProtocol = process.env.NODE_ENV === "production" ? "wss://" : "ws://";
export const siteName = process.env.NEXT_PUBLIC_API_URL;
export const socketEndpoint = process.env.NEXT_PUBLIC_API_URL.toString().replace(/.+\/{2}/, wsProtocol).concat("/ws");
export const algodProtocol = process.env.NEXT_PUBLIC_ALGOD_PROTOCOL;
export const algodAddr = process.env.NEXT_PUBLIC_ALGOD_ADDR;
export const algodToken = process.env.NEXT_PUBLIC_ALGOD_TOKEN;
export const isLocal = algodAddr.indexOf("0.0.0.0") > -1 || algodAddr.indexOf("127.0.0.1") > -1 || algodAddr.indexOf("localhost") > -1;
