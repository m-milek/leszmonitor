import { expect } from "@playwright/test";
import test from "../fixtures/leszmonitorFixture";
import { loginAs } from "../helpers/auth";
import { createMonitor, pushUrl, uniqueName } from "../helpers/monitors";

const PUBLIC_URL = "https://status.example.com";

test.describe("Instance Settings", () => {
  test("Push URL uses the configured public URL", async ({
    page,
    adminAuth,
  }) => {
    const jwt = await loginAs(page, adminAuth);

    try {
      await page.goto("/admin");
      await page.getByLabel("Public URL").fill(PUBLIC_URL);
      await page.getByRole("button", { name: "Save" }).click();
      await expect(page.getByText("Global parameters saved")).toBeVisible();

      await createMonitor(page, {
        type: "Push",
        name: uniqueName("Public URL Push Monitor"),
      });

      await expect(pushUrl(page)).toContainText(`${PUBLIC_URL}/api/v1/push/`);
    } finally {
      await page.request.patch("/api/v1/global-parameters", {
        headers: { Authorization: `Bearer ${jwt}` },
        data: { "instance.public_url": null },
      });
    }
  });
});
