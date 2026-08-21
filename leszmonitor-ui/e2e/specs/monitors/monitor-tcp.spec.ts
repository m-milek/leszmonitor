import { expect } from "@playwright/test";
import test from "../../fixtures/leszmonitorFixture.ts";

test.describe("TCP Monitor", () => {
  test("Creates a valid TCP monitor", async ({ page }) => {
    await page.goto("/projects/leszmaks-sandbox/monitors/new");

    const randomMonitorNumber = Math.floor(Math.random() * 10000);

    await page
      .getByRole("combobox")
      .filter({ hasText: "Select Monitor Type" })
      .click();
    await page.getByRole("option", { name: "TCP" }).click();
    await page
      .getByLabel("Name")
      .fill(`Test TCP Monitor ${randomMonitorNumber}`);
    await page.getByLabel("Host").fill("example.com");
    await page.getByLabel("Port").fill("80");

    await page.getByText("Create Monitor").click();

    await expect(page).toHaveURL(
      /\/projects\/leszmaks-sandbox\/monitors\/test-tcp-monitor-\d+$/,
    );
  });
});
