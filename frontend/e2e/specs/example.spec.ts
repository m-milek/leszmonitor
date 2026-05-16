import { expect } from "@playwright/test";
import test from "../fixtures/leszmonitorFixture";

test("has title", async ({ page }) => {
  await page.goto("http://localhost:7001/login");

  await expect(page).toHaveTitle(/Leszmonitor/);
});
