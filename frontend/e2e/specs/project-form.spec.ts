import {expect} from "@playwright/test";
import test from "../fixtures/leszmonitorFixture.ts";

test.describe("Project Form", () => {
  test("Successfully creates a new project via happy path", async ({page}) => {
    await page.goto("/projects");

    await page.getByRole("button", {name: "Create Project"}).click();

    const projectName = `E2E Test Project ${Date.now()}`;
    await page.getByLabel("Project Name").fill(projectName);
    await page.getByLabel("Description (Optional)").fill("A test project description for E2E purposes.");

    await page.getByRole("dialog").getByRole("button", {name: "Create Project"}).click();

    await expect(page.getByText(projectName)).toBeVisible();
  });

  test("Shows validation errors when submitting an empty form", async ({page}) => {
    await page.goto("/projects");

    await page.getByRole("button", {name: "Create Project"}).click();

    await page.getByRole("dialog").getByRole("button", {name: "Create Project"}).click();

    await expect(page.getByText("Project name is required")).toBeVisible();
    await expect(page.getByText("Slug is required")).toBeVisible();
  });

  test.skip("Shows an error message when trying to create a project that already exists", async ({page}) => {
    await page.goto("/projects");

    await page.getByRole("button", {name: "Create Project"}).click();
    const projectName = `Duplicate Test Project ${Date.now()}`;
    await page.getByLabel("Project Name").fill(projectName);
    await page.getByRole("dialog").getByRole("button", {name: "Create Project"}).click();

    await expect(page.getByText(projectName)).toBeVisible();

    await page.getByRole("button", {name: "Create Project"}).click();
    await page.getByLabel("Project Name").fill(projectName);
    await page.getByRole("dialog").getByRole("button", {name: "Create Project"}).click();

    /* Check for error message when the error handling is ready*/
  });
});
