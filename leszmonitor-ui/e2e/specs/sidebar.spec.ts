import { expect } from "@playwright/test";
import test from "../fixtures/leszmonitorFixture";

test.describe("Sidebar", () => {
  test("Home icon navigates to the monitors list", async ({ page }) => {
    await page.goto("/monitors/new");

    await page.getByRole("link", { name: "Home" }).click();

    await expect(page.getByText("Monitors")).toBeVisible();
    expect(page.url()).toMatch(/\/monitors$/);
  });
});
