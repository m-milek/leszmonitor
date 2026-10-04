import { expect, type Page } from "@playwright/test";

export const openNewTagDialog = async (page: Page) => {
  await page.goto("/tags");
  await page.getByRole("button", { name: "Add Tag" }).click();

  return page.getByRole("dialog", { name: "Add New Tag" });
};

export const createTag = async (page: Page, name: string) => {
  const dialog = await openNewTagDialog(page);
  await dialog.getByLabel("Name").fill(name);
  await dialog.getByRole("button", { name: "Create Tag" }).click();

  await expect(dialog).toBeHidden();
};

export const selectTagInMonitorForm = async (page: Page, name: string) => {
  await page.getByLabel("Tags").click();
  await page.getByRole("option", { name }).click();
  await page.keyboard.press("Escape");
};
