import { type PlaywrightTestConfig } from "@playwright/test";

const config: PlaywrightTestConfig = {
  use: {
    baseURL: "https://plateauview.mlit.go.jp",
    screenshot: "only-on-failure",
    video: "retain-on-failure",
  },
  testDir: "e2e",
  globalSetup: "./e2e/utils/setup.ts",
  reporter: process.env.CI ? "github" : "list",
};

export default config;
