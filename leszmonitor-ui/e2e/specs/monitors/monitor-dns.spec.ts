import { expect } from "@playwright/test";
import test from "../../fixtures/leszmonitorFixture.ts";

test.describe("DNS Monitor", () => {
  test("Creates a valid DNS monitor", async ({ page }) => {
    await page.goto("/projects/leszmaks-sandbox/monitors/new");

    const randomMonitorNumber = Math.floor(Math.random() * 10000);

    await page
      .getByRole("combobox")
      .filter({ hasText: "Select Monitor Type" })
      .click();
    await page.getByRole("option", { name: "DNS" }).click();
    await page
      .getByLabel("Name", { exact: true })
      .fill(`Test DNS Monitor ${randomMonitorNumber}`);
    await page.getByLabel("Hostname").fill("example.com");

    await page.getByText("Create Monitor").click();

    await expect(page).toHaveURL(
      /\/projects\/leszmaks-sandbox\/monitors\/test-dns-monitor-\d+$/,
    );
  });
});
