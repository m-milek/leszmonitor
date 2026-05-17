import fs from "fs";
import path from "path";

const AUTH_USERNAME = process.env.AUTH_USERNAME || "";
const AUTH_PASSWORD = process.env.AUTH_PASSWORD || "";

const BASE_URL = "http://localhost:7001";
const MOCKS_BASE_PATH = "../../frontend/e2e/mocks"

interface MockMapping {
  url: string;
  params?: Record<string, string>;
  output: string;
}

const mappings: MockMapping[] = [
  {url: "/api/v1/users", output: `${MOCKS_BASE_PATH}/users.json`},
  {url: "/api/v1/users/leszmak", output: `${MOCKS_BASE_PATH}/user_leszmak.json`},
  {url: "/api/v1/projects", output: `${MOCKS_BASE_PATH}/projects.json`},
  {
    url: "/api/v1/monitors",
    params: {projectSlug: "leszmaks-sandbox"},
    output: `${MOCKS_BASE_PATH}/leszmaks_sandbox_monitors.json`
  },
  {
    url: "/api/v1/projects/leszmaks-sandbox/monitors/gnu",
    output: `${MOCKS_BASE_PATH}/leszmaks_sandbox_monitor_gnu.json`
  },
  // {
  //   url: "/api/v1/monitors/c71649c4-cbb6-4f1d-8149-210a365f2a99/results",
  //   params: {page: "1", per_page: "100"},
  //   output: `${MOCKS_BASE_PATH}/leszmaks_sandbox_monitor_gnu_100_results.json`
  // },
];

async function getLoginToken(): Promise<string> {
  if (!AUTH_USERNAME || !AUTH_PASSWORD) {
    throw new Error("AUTH_USERNAME and AUTH_PASSWORD environment variables must be set");
  }

  console.log(`Username: ${AUTH_USERNAME}\nPassword: ${AUTH_PASSWORD}`);

  const url = `${BASE_URL}/api/v1/auth/login`;
  const res = await fetch(url, {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({username: AUTH_USERNAME, password: AUTH_PASSWORD}),
  });

  if (!res.ok) {
    throw new Error(`Login failed: ${res.status} ${res.statusText}`);
  }

  const data = await res.json();

  return data.jwt;
}

async function fetchAndSave(m: MockMapping, loginToken: string): Promise<void> {
  const url = `${BASE_URL}${m.url}`;
  const u = new URL(url);
  if (m.params) {
    Object.entries(m.params).forEach(([k, v]) => u.searchParams.set(k, v));
  }

  console.log(`GET ${u} -> ${m.output}`);

  const headers: Record<string, string> = {Accept: "application/json"};

  headers["Authorization"] = `Bearer ${loginToken}`;

  try {
    const res = await fetch(u.toString(), {method: "GET", headers});

    if (!res.ok) {
      console.log(`  FAIL ${res.url}: ${res.status} ${res.statusText}`);
      return;
    }

    const data = await res.json();
    const pretty = JSON.stringify(data, null, 2);

    fs.mkdirSync(path.dirname(m.output), {recursive: true});
    fs.writeFileSync(m.output, pretty, "utf-8");

    console.log("  OK");
  } catch (err) {
    console.log(`  FAIL ${(err as Error).message}`);
  }
}

async function main(): Promise<void> {
  console.log("=== Mock Sync ===\n");

  const loginToken = await getLoginToken();

  for (const m of mappings) {
    await fetchAndSave(m, loginToken);
  }

  console.log("\n=== Done ===");
}

main();