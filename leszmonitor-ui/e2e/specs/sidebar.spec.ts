import { expect } from "@playwright/test";
import test from "../fixtures/leszmonitorFixture";

test.describe("Sidebar", () => {
  test("Monitors link navigates to the monitors list", async ({ page }) => {
    await page.goto("/monitors/new");

    await page.getByRole("link", { name: "Monitors" }).click();

    await expect(page.getByRole("heading", { name: "Monitors" })).toBeVisible();
    expect(page.url()).toMatch(/\/monitors$/);
  });
});
