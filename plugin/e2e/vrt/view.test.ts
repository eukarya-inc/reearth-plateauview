import { test, expect } from "@playwright/test";

test("vrt: top page", async ({ page }) => {
  await page.goto("https://plateauview.mlit.go.jp/");
  // close button
  // await page.click("svg");
  await page.waitForLoadState("networkidle");
  // await page.waitForSelector("text=カタログから検索する");
  const image = await page.screenshot();
  expect(image).toMatchSnapshot();
});
