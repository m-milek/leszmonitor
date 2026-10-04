import { expect } from "@playwright/test";
import test from "../fixtures/leszmonitorFixture";
import { createMonitor, uniqueName } from "../helpers/monitors";
import {
  createTag,
  openNewTagDialog,
  selectTagInMonitorForm,
} from "../helpers/tags";

test.describe("Tags", () => {
  test("Creates a tag with a live preview", async ({ page }) => {
    const tagName = uniqueName("e2e tag");

    const dialog = await openNewTagDialog(page);
    await dialog.getByLabel("Name").fill(tagName);
    await expect(dialog.getByText(tagName)).toBeVisible();

    await dialog.getByRole("button", { name: "Create Tag" }).click();

    await expect(dialog).toBeHidden();
    await expect(page.getByRole("table").getByText(tagName)).toBeVisible();
  });

  test("Filters monitors by tag", async ({ page }) => {
    const tagName = uniqueName("e2e filter tag");
    const taggedName = uniqueName("Tagged Push Monitor");
    const untaggedName = uniqueName("Untagged Push Monitor");

    await createTag(page, tagName);
    await createMonitor(page, {
      type: "Push",
      name: taggedName,
      fillFields: () => selectTagInMonitorForm(page, tagName),
    });
    await createMonitor(page, { type: "Push", name: untaggedName });

    await page.goto("/monitors");
    await expect(page.getByRole("link", { name: untaggedName })).toBeVisible();

    await page.getByRole("combobox").filter({ hasText: "Tags" }).click();
    await page.getByRole("option", { name: tagName }).click();
    await page.keyboard.press("Escape");

    await expect(page.getByRole("link", { name: taggedName })).toBeVisible();
    await expect(page.getByRole("link", { name: untaggedName })).toBeHidden();
  });
});
