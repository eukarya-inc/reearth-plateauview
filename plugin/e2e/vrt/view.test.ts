import { test, expect } from "@playwright/test";

test("should match previous screenshot", async ({ page }) => {
  await page.goto("/");
  const image = await page.screenshot();
  expect(image).toMatchSnapshot();
});
