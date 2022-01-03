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
export const siteName = process.env.NEXT_PUBLIC_API_URL;
