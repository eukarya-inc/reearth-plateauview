import { devices, type PlaywrightTestConfig } from "@playwright/test";

const config: PlaywrightTestConfig = {
  use: {
    baseURL: "https://plateauview.mlit.go.jp",
    screenshot: "only-on-failure",
    video: "retain-on-failure",
  },
  testDir: "e2e",
  globalSetup: "./e2e/utils/setup.ts",
  reporter: process.env.CI ? "github" : "list",
  projects: [
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
        headless: false,
        launchOptions: {
          args: ["--no-sandbox"],
        },
      },
    },
  ],
};

export default config;
