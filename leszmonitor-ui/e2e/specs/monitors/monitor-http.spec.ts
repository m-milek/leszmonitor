import { expect } from "@playwright/test";
import test from "../../fixtures/leszmonitorFixture";

test.describe("HTTP Monitor", () => {
  test("Creates a valid HTTP monitor", async ({ page }) => {
    await page.goto("/monitors/new");

    const randomMonitorNumber = Math.floor(Math.random() * 10000);

    await page
      .getByRole("combobox")
      .filter({ hasText: "Select Monitor Type" })
      .click();
    await page.getByRole("option", { name: "HTTP" }).click();
    await page
      .getByLabel("Name")
      .fill(`Test HTTP Monitor ${randomMonitorNumber}`);
    await page.getByLabel("URL").fill("https://example.com");
    await page.getByRole("combobox", { name: "Method" }).click();
    await page.getByRole("option", { name: "GET" }).click();

    await page.getByLabel("Expected Status Codes").click();
    await page.getByRole("option", { name: "200" }).click();

    await page.getByText("Create Monitor").click();

    await expect(page).toHaveURL(/\/monitors\/test-http-monitor-\d+$/);
  });
});
