import { describe, expect, test } from "vitest";
import { prepareResults } from "./prepare-results.ts";
import type { MonitorResult } from "@/lib/types.ts";

describe("prepareResults", () => {
  const createResult = (dateStr: string): MonitorResult =>
    ({
      id: "test-id",
      monitorId: "test-monitor-id",
      status: "up",
      isManuallyTriggered: false,
      durationMs: 100,
      errorDetails: { errorMessage: "", errors: [], failures: [], message: "" },
      monitorType: "http",
      details: { statusCode: 200, contentLength: 10, proto: "HTTP/1.1" },
      createdAt: new Date(dateStr),
    }) as MonitorResult;

  test("returns empty array when length is 0", () => {
    expect(prepareResults([], 0)).toEqual([]);
  });

  test("returns array with undefined when length is greater than results", () => {
    const results = [createResult("2024-01-01T10:00:00Z")];
    const res = prepareResults(results, 3);

    expect(res).toHaveLength(3);
    expect(res[0]).toBeUndefined();
    expect(res[1]).toBeUndefined();
    expect(res[2]?.createdAt.getTime()).toBe(
      new Date("2024-01-01T10:00:00Z").getTime(),
    );
  });

  test("sorts and returns most recent items in ascending order", () => {
    const results = [
      createResult("2024-01-01T10:00:00Z"),
      createResult("2024-01-01T12:00:00Z"),
      createResult("2024-01-01T11:00:00Z"),
    ];

    // We want 2 items, which should be the most recent ones: 12:00 and 11:00.
    // They should be returned in ascending order: 11:00 then 12:00.
    const res = prepareResults(results, 2);
    expect(res).toHaveLength(2);
    expect(res[0]?.createdAt.getTime()).toBe(
      new Date("2024-01-01T11:00:00Z").getTime(),
    );
    expect(res[1]?.createdAt.getTime()).toBe(
      new Date("2024-01-01T12:00:00Z").getTime(),
    );
  });

  test("returns array of exact length when length matches results", () => {
    const results = [
      createResult("2024-01-01T10:00:00Z"),
      createResult("2024-01-01T11:00:00Z"),
    ];

    const res = prepareResults(results, 2);
    expect(res).toHaveLength(2);
    expect(res[0]?.createdAt.getTime()).toBe(
      new Date("2024-01-01T10:00:00Z").getTime(),
    );
    expect(res[1]?.createdAt.getTime()).toBe(
      new Date("2024-01-01T11:00:00Z").getTime(),
    );
  });

  test("handles negative length by returning empty array", () => {
    const results = [createResult("2024-01-01T10:00:00Z")];
    expect(prepareResults(results, -5)).toEqual([]);
  });

  test("handles float length by flooring it", () => {
    const results = [
      createResult("2024-01-01T10:00:00Z"),
      createResult("2024-01-01T11:00:00Z"),
      createResult("2024-01-01T12:00:00Z"),
    ];
    const res = prepareResults(results, 2.7);
    expect(res).toHaveLength(2);
    expect(res[0]?.createdAt.getTime()).toBe(
      new Date("2024-01-01T11:00:00Z").getTime(),
    );
    expect(res[1]?.createdAt.getTime()).toBe(
      new Date("2024-01-01T12:00:00Z").getTime(),
    );
  });
});
