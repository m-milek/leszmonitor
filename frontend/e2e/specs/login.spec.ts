import { expect } from "@playwright/test";
import test from "../fixtures/leszmonitorFixture";

test.describe("Login", () => {
  test("Displays app title", async ({ page }) => {
    await page.goto("/login");

    await expect(page).toHaveTitle(/Leszmonitor/);
  });

  test("Good credentials allow to log in", async ({ page, auth }) => {
    await page.goto("/login");

    await page.getByLabel("Username").fill(auth.username);
    await page.getByLabel("Password").fill(auth.password);
    await page.getByRole("button", { name: "Log in" }).click();

    // Redirects to projects right away
    await expect(page).toHaveURL("/projects");
  });

  test("Bad credentials don't log in", async ({ page }) => {
    await page.goto("/login");

    await page.getByLabel("Username").fill("wrong");
    await page.getByLabel("Password").fill("credentials");
    await page.getByRole("button", { name: "Log in" }).click();

    await expect(page).toHaveURL("/login");

    await expect(
      page.getByText(
        "Failed to log in. Please check your credentials and try again.",
      ),
    ).toBeVisible();
  });

  test("Empty credentials show validation errors", async ({ page }) => {
    await page.goto("/login");

    await page.getByRole("button", { name: "Log in" }).click();

    await expect(page.getByText("Username is required")).toBeVisible();
    await expect(page.getByText("Password is required")).toBeVisible();
  });
});
