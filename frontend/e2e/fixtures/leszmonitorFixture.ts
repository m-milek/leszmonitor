import { test as base } from "@playwright/test";

import users from "../mocks/users.json" with { type: "json" };
import projects from "../mocks/projects.json" with { type: "json" };
import userLeszmak from "../mocks/user_leszmak.json" with { type: "json" };
import leszmaksSandboxMonitors from "../mocks/leszmaks_sandbox_monitors.json" with { type: "json" };
import leszmaksSandboxMonitorGnu from "../mocks/leszmaks_sandbox_monitor_gnu.json" with { type: "json" };
import leszmaksSandboxMonitorGnu100Results from "../mocks/leszmaks_sandbox_monitor_gnu_100_results.json" with { type: "json" };

type LeszmonitorFixture = {
  dupa: string;
};

export const test = base.extend<LeszmonitorFixture>({
  page: async ({ page }, use) => {
    await page.route("**/api/v1/users", (route) => {
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(users),
      });
    });
    await page.route("**/api/v1/users/leszmak", (route) => {
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(userLeszmak),
      });
    });
    await page.route("**/api/v1/projects", (route) => {
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify(projects),
      });
    });
    await page.route(
      "**/api/v1/monitors?projectSlug=leszmaks-sandbox",
      (route) => {
        route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(leszmaksSandboxMonitors),
        });
      },
    );
    await page.route(
      "**/api/v1/projects/leszmaks-sandbox/monitors/gnu",
      (route) => {
        route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(leszmaksSandboxMonitorGnu),
        });
      },
    );
    await page.route(
      "**/api/v1/monitors/c71649c4-cbb6-4f1d-8149-210a365f2a99/results?page=1&per_page=100",
      (route) => {
        route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify(leszmaksSandboxMonitorGnu100Results),
        });
      },
    );

    await use(page);
  },
});

export { test as default };
