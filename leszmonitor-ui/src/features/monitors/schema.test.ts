import { describe, expect, it } from "vitest";
import {
  defaultConfigs,
  mapMonitorToFormValues,
  newMonitorSchema,
  newMonitorSchemaDefaultValues
} from "@/features/monitors/schema";
import type { Monitor, MonitorType } from "@/features/monitors/types";

const serverMonitor = (overrides: Partial<Monitor>): Monitor =>
  ({
    id: "8f1c2d3e",
    ownerId: "owner-1",
    name: "Example",
    slug: "example",
    interval: 60,
    runState: "active",
    createdAt: new Date("2026-01-01T00:00:00Z"),
    updatedAt: new Date("2026-01-01T00:00:00Z"),
    ...overrides,
  }) as Monitor;

const serverMonitors: Record<MonitorType, Monitor> = {
  http: serverMonitor({
    type: "http",
    probeConfig: {
      method: "GET",
      url: "https://example.com",
      saveResponseBody: false,
      saveResponseHeaders: false,
    },
  }),
  tcp: serverMonitor({
    type: "tcp",
    probeConfig: {
      host: "example.com",
      port: 8080,
      protocol: "tcp",
      timeout: 5000,
      retryCount: 1,
    },
  }),
  dns: serverMonitor({
    type: "dns",
    probeConfig: {
      hostname: "example.com",
      dnsServer: "1.1.1.1",
      recordType: "A",
      expectedRecordValues: [],
    },
  }),
};

const issues = (input: unknown) => {
  const result = newMonitorSchema.safeParse(input);
  return result.success ? [] : result.error.issues.map((i) => i.message);
};

const baseFormValues = {
  ...newMonitorSchemaDefaultValues,
  name: "Example",
  slug: "example",
};

const tcpForm = (probe: Record<string, unknown> = {}) => ({
  ...baseFormValues,
  type: "tcp",
  probeConfig: {
    host: "example.com",
    port: 443,
    protocol: "tcp",
    timeout: 5000,
    retryCount: 0,
    ...probe,
  },
});

describe("mapMonitorToFormValues", () => {
  it.each(["http", "tcp", "dns"] as const)(
    "maps a %s monitor from the server into values the schema accepts",
    (type) => {
      const values = mapMonitorToFormValues(serverMonitors[type]);
      expect(issues(values)).toEqual([]);
      expect(values.type).toBe(type);
    },
  );

  it("keeps server values and backfills the rest from the type defaults", () => {
    const partial = serverMonitor({
      type: "tcp",
      probeConfig: { host: "db.internal", port: 5432 },
    } as Partial<Monitor>);

    expect(mapMonitorToFormValues(partial).probeConfig).toEqual({
      host: "db.internal",
      port: 5432,
      protocol: "tcp",
      timeout: 5000,
      retryCount: 3,
    });
  });
});

describe("defaultConfigs", () => {
  it.each([
    ["http", { url: "https://example.com" }],
    ["tcp", { host: "example.com" }],
    ["dns", { hostname: "example.com" }],
  ] as const)(
    "the %s defaults only miss the required field",
    (type, filled) => {
      const values = {
        ...baseFormValues,
        type,
        probeConfig: { ...defaultConfigs[type], ...filled },
      };

      expect(issues(values)).toEqual([]);
    },
  );
});

describe("newMonitorSchema", () => {
  it.each([
    [1, true],
    [65535, true],
    [0, false],
    [65536, false],
  ])("accepts port %i: %s", (port, valid) => {
    expect(newMonitorSchema.safeParse(tcpForm({ port })).success).toBe(valid);
  });

  it.each(["Example", "with space", "trailing-", "-leading", "under_score"])(
    "rejects the slug %j",
    (slug) => {
      expect(issues({ ...tcpForm(), slug })).toContain(
        "Invalid slug format. Must be lowercase, alphanumeric, and can include hyphens.",
      );
    },
  );

  it("rejects an interval below one second", () => {
    expect(issues({ ...tcpForm(), interval: 0 })).toContain(
      "Interval must be at least 1 second",
    );
  });

  it("rejects a malformed url on an http monitor", () => {
    const values = {
      ...baseFormValues,
      type: "http",
      probeConfig: {
        ...defaultConfigs.http,
        url: "not-a-url",
      },
    };

    expect(issues(values)).toContain("Invalid URL");
  });

  it("rejects a probe config belonging to a different monitor type", () => {
    const mismatched = {
      ...baseFormValues,
      type: "tcp",
      probeConfig: { ...defaultConfigs.http, url: "https://example.com" },
    };

    expect(newMonitorSchema.safeParse(mismatched).success).toBe(false);
  });
});
