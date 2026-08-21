import test from "../fixtures/leszmonitorFixture";
import { expect } from "@playwright/test";

test.describe("Register", () => {
  test("Registers new user", async ({ page }) => {
    await page.goto("/register");

    const username = `testuser${Date.now()}`;
    const password = "TestPassword123!";

    await page.getByLabel("Username").fill(username);
    await page.locator("#password").fill(password);
    await page.locator("#passwordConfirm").fill(password);
    await page.getByRole("button", { name: "Register" }).click();

    await expect(page).toHaveURL("/projects");
  });

  test("Fails to register with a short password", async ({ page }) => {
    await page.goto("/register");

    const username = `testuser${Date.now()}`;
    const password = "short";

    await page.getByLabel("Username").fill(username);
    await page.locator("#password").fill(password);
    await page.locator("#passwordConfirm").fill(password);
    await page.getByRole("button", { name: "Register" }).click();

    await expect(
      page.getByText("Password has to be at least 6 characters long"),
    ).toBeVisible();
  });

  test("Fails to register with mismatched passwords", async ({ page }) => {
    await page.goto("/register");

    const username = `testuser${Date.now()}`;
    const password = "TestPassword123!";
    const passwordConfirm = "DifferentPassword123!";

    await page.getByRole("textbox", { name: "Username" }).fill(username);
    await page.locator("#password").fill(password);
    await page.locator("#passwordConfirm").fill(passwordConfirm);
    await page.getByRole("button", { name: "Register" }).click();

    await expect(page.getByText("Passwords don't match")).toBeVisible();
  });

  test("Fails to register with an existing username", async ({
    page,
    auth,
  }) => {
    await page.goto("/register");

    const username = auth.username;
    const password = auth.password;

    await page.getByLabel("Username").fill(username);
    await page.locator("#password").fill(password);
    await page.locator("#passwordConfirm").fill(password);
    await page.getByRole("button", { name: "Register" }).click();

    await expect(
      page.getByText("Registration failed. Please try again."),
    ).toBeVisible();
  });
});
