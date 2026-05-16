import { expect } from "@playwright/test";
import test from "../fixtures/leszmonitorFixture";

test.describe("Login", () => {
  test("Displays app title", async ({ page }) => {
    await page.goto("/login");

    await expect(page).toHaveTitle(/Leszmonitor/);
  });
});
