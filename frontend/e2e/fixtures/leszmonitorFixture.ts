import { test as base } from "@playwright/test";

import users from "../mocks/users.json" with { type: "json" };
import projects from "../mocks/projects.json" with { type: "json" };
import userLeszmak from "../mocks/user_leszmak.json" with { type: "json" };
import leszmaksSandboxMonitors from "../mocks/leszmaks_sandbox_monitors.json" with { type: "json" };
import leszmaksSandboxMonitorGnu from "../mocks/leszmaks_sandbox_monitor_gnu.json" with { type: "json" };
import leszmaksSandboxMonitorGnu100Results from "../mocks/leszmaks_sandbox_monitor_gnu_100_results.json" with { type: "json" };

interface AuthConfig {
  username: string;
  password: string;
}

type LeszmonitorFixture = {
  auth: AuthConfig;
};

export const test = base.extend<LeszmonitorFixture>({
  auth: {
    username: "leszmak",
    password: "123123",
  },
  page: async ({ page }, use) => {
    await page.route("**/api/v1/auth/login", (route) => {
      if (route.request().method() === "POST") {
        const postData = route.request().postData();
        if (postData) {
          const { username, password } = JSON.parse(postData);
          if (username !== "leszmak" || password !== "123123") {
            route.fulfill({
              status: 401,
              contentType: "application/json",
              body: JSON.stringify({ error: "Invalid credentials" }),
            });
            return;
          }
        } else {
          route.fulfill({
            status: 400,
            contentType: "application/json",
            body: JSON.stringify({ error: "Missing request body" }),
          });
          return;
        }
      }
      route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          jwt: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzkwNDc1MzgsImlhdCI6MTc3ODk2MTEzOCwidXNlcm5hbWUiOiJsZXN6bWFrIn0.BD0zIfx2kdtPd_GZwsVYh30HxdQf4sm4JYGFFNttT0M",
        }),
      });
    });
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
