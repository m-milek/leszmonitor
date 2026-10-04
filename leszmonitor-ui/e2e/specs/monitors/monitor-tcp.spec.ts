import test from "../../fixtures/leszmonitorFixture";
import { createMonitor, uniqueName } from "../../helpers/monitors";

test.describe("TCP Monitor", () => {
  test("Creates a valid TCP monitor", async ({ page }) => {
    await createMonitor(page, {
      type: "TCP",
      name: uniqueName("Test TCP Monitor"),
      fillFields: async () => {
        await page.getByLabel("Host").fill("example.com");
        await page.getByLabel("Port").fill("80");
      },
    });
  });
});
