import { expect } from "@playwright/test";
import test from "../fixtures/leszmonitorFixture";

test.describe("Logout", () => {
  test("Logging out drops the session and blocks protected routes", async ({
    page,
  }) => {
    await page.goto("/monitors");

    await page.getByRole("button", { name: "User menu" }).click();
    await page.getByRole("menuitem", { name: "Log out" }).click();

    await expect(page).toHaveURL("/login");

    const cookies = await page.context().cookies();
    expect(cookies.find((c) => c.name === "LOGIN_TOKEN")).toBeUndefined();

    await page.goto("/monitors");
    await expect(page).toHaveURL("/login");
  });
});
