import test from "../../fixtures/leszmonitorFixture";
import { createMonitor, uniqueName } from "../../helpers/monitors";

test.describe("DNS Monitor", () => {
  test("Creates a valid DNS monitor", async ({ page }) => {
    await createMonitor(page, {
      type: "DNS",
      name: uniqueName("Test DNS Monitor"),
      fillFields: () => page.getByLabel("Hostname").fill("example.com"),
    });
  });
});
