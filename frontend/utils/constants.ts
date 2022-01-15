declare global {
  namespace NodeJS {
    interface ProcessEnv {
      // @ts-ignore
      NODE_ENV: "development" | "production" | "test";
      NEXT_PUBLIC_API_URL: string;
    }
  }
}

// @ts-ignore
const wsProtocol = process.env.NODE_ENV === "production" ? "wss://" : "ws://";
export const siteName = process.env.NEXT_PUBLIC_API_URL;
export const socketEndpoint = process.env.NEXT_PUBLIC_API_URL.toString().replace(/.+\/{2}/, wsProtocol).concat("/ws");
