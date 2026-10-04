import test from "../../fixtures/leszmonitorFixture";
import { createMonitor, uniqueName } from "../../helpers/monitors";

test.describe("HTTP Monitor", () => {
  test("Creates a valid HTTP monitor", async ({ page }) => {
    await createMonitor(page, {
      type: "HTTP",
      name: uniqueName("Test HTTP Monitor"),
      fillFields: async () => {
        await page.getByLabel("URL").fill("https://example.com");
        await page.getByRole("combobox", { name: "Method" }).click();
        await page.getByRole("option", { name: "GET" }).click();

        await page.getByLabel("Expected Status Codes").click();
        await page.getByRole("option", { name: "200" }).click();
      },
    });
  });
});
