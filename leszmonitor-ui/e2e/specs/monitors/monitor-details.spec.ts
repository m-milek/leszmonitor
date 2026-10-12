import { expect } from "@playwright/test";
import test from "../../fixtures/leszmonitorFixture";
import { createMonitor, uniqueName } from "../../helpers/monitors";

test.describe("Monitor Details", () => {
  test("Keeps the selected results range in the URL", async ({ page }) => {
    await createMonitor(page, {
      type: "Push",
      name: uniqueName("Range Push Monitor"),
    });
    await expect(page).toHaveURL(/[?&]range=86400/);

    await page.getByRole("combobox").filter({ hasText: "Last 1 day" }).click();
    await page.getByRole("option", { name: "Last 1 hour" }).click();
    await expect(page).toHaveURL(/[?&]range=3600/);

    await page.reload();

    await expect(page).toHaveURL(/[?&]range=3600/);
    await expect(
      page.getByRole("combobox").filter({ hasText: "Last 1 hour" }),
    ).toBeVisible();
  });

  test("Exports the monitor as YAML", async ({ page }) => {
    const name = uniqueName("Export Push Monitor");
    await createMonitor(page, {
      type: "Push",
      name,
      fillFields: () => page.getByLabel("Grace period (s)").fill("60"),
    });

    await page.getByRole("button", { name: "More actions" }).click();
    await page.getByRole("menuitem", { name: "Export" }).click();

    const dialog = page.getByRole("dialog");
    await expect(dialog).toContainText(`name: ${name}`);
    await expect(dialog).toContainText("type: push");
    await expect(dialog).toContainText("gracePeriodSeconds: 60");
  });
});
