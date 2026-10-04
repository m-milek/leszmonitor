import { expect, type Page } from "@playwright/test";

export type MonitorTypeOption = "HTTP" | "TCP" | "DNS" | "Push";

export interface CreateMonitorOptions {
  type: MonitorTypeOption;
  name: string;
  fillFields?: () => Promise<void>;
}

export const uniqueName = (prefix: string) => `${prefix} ${Date.now()}`;

const slugFromName = (name: string) => name.toLowerCase().replaceAll(" ", "-");

export const createMonitor = async (
  page: Page,
  { type, name, fillFields }: CreateMonitorOptions,
) => {
  await page.goto("/monitors/new");

  await page
    .getByRole("combobox")
    .filter({ hasText: "Select Monitor Type" })
    .click();
  await page.getByRole("option", { name: type }).click();
  await page.getByLabel("Name", { exact: true }).fill(name);
  await fillFields?.();

  await page.getByText("Create Monitor").click();

  await expect(page).toHaveURL(
    new RegExp(`/monitors/${slugFromName(name)}(\\?|$)`),
  );
};

export const pushUrl = (page: Page) =>
  page.locator("code", { hasText: "/api/v1/push/" });

export const getPushMonitorId = async (page: Page) => {
  const url = await pushUrl(page).innerText();
  return url.slice(url.lastIndexOf("/") + 1);
};
