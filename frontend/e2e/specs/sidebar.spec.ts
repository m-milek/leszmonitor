import { expect } from "@playwright/test";
import test from "../fixtures/leszmonitorFixture";

test.describe("Sidebar", () => {
  test("Project selector selects a project", async ({ page }) => {
    await page.goto("/projects");

    await page.getByRole("combobox").click();
    await page.getByRole("option", { name: "leszmak's Sandbox" }).click();

    await expect(
      page.getByLabel("leszmak's Sandbox").getByText("leszmak's Sandbox"),
    ).toBeVisible();
    expect(page.url()).toMatch(/\/projects\/leszmaks-sandbox$/);
  });

  test("Home icon navigates to project list", async ({ page }) => {
    await page.goto("/projects/leszmaks-sandbox");

    await page.getByRole("button", { name: "Home" }).click();

    await expect(page.getByText("Your Projects")).toBeVisible();
    expect(page.url()).toMatch(/\/projects$/);
  });
});
