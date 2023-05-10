import { test, expect } from "@playwright/test";

test("vrt: top page", async ({ page }) => {
  await page.goto("");
  const image = await page.screenshot();
  expect(image).toMatchSnapshot();
});
