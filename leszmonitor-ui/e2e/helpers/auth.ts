import { expect, test, type Page } from "@playwright/test";
import type { AuthConfig } from "../fixtures/leszmonitorFixture";

export const loginAs = async (page: Page, auth: AuthConfig) => {
  const res = await page.request.post("/api/v1/auth/login", { data: auth });
  expect(res.ok()).toBe(true);
  const { jwt } = (await res.json()) as { jwt: string };

  await page.context().addCookies([
    {
      url: test.info().project.use.baseURL!,
      name: "LOGIN_TOKEN",
      value: jwt,
    },
  ]);

  return jwt;
};
