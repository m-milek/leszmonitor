import { expect } from "@playwright/test";
import test from "../../fixtures/leszmonitorFixture";
import {
  createMonitor,
  getPushMonitorId,
  pushUrl,
  uniqueName,
} from "../../helpers/monitors";

test.describe("Push Monitor", () => {
  test("Creates a valid push monitor", async ({ page }) => {
    await createMonitor(page, {
      type: "Push",
      name: uniqueName("Test Push Monitor"),
      fillFields: () => page.getByLabel("Grace period (s)").fill("60"),
    });

    await expect(pushUrl(page)).toBeVisible();
  });

  test("Shows a failure reported through the push URL", async ({ page }) => {
    await createMonitor(page, {
      type: "Push",
      name: uniqueName("Reporting Push Monitor"),
    });
    const monitorId = await getPushMonitorId(page);

    const res = await page.request.post(
      `/api/v1/push/${monitorId}?status=down`,
      { data: "backup failed" },
    );
    expect(res.status()).toBe(202);

    await expect(
      page.getByRole("cell", { name: "Self-reported failure" }),
    ).toBeVisible();
  });

  test("Rejects a push to a paused monitor", async ({ page }) => {
    await createMonitor(page, {
      type: "Push",
      name: uniqueName("Paused Push Monitor"),
    });
    const monitorId = await getPushMonitorId(page);

    await page.getByRole("button", { name: "Pause", exact: true }).click();
    await expect(
      page.getByRole("button", { name: "Resume", exact: true }),
    ).toBeVisible();

    const res = await page.request.get(`/api/v1/push/${monitorId}`);
    expect(res.status()).toBe(409);
  });
});
