import { test as base } from "@playwright/test";

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
    await page.context().addCookies([
      {
        url: "http://localhost:3000",
        name: "LOGIN_TOKEN",
        value:
          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUzNzkxMzc2NDksImlhdCI6MTc3OTE0MTI0OSwidXNlcm5hbWUiOiJsZXN6bWFrIn0.f4hIqBGicZyvJ5UYIQ9-33GsndjA_mKMLXN5WBHHUvg",
      },
    ]);
    await use(page);
  },
});

export { test as default };
