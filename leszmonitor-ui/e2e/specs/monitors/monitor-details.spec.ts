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
});
