import type { MonitorResult } from "@/lib/types.ts";

export const prepareResults = (
  monitorResults: MonitorResult[],
  length: number,
): Array<MonitorResult | undefined> => {
  const safeLength = Math.max(0, Math.floor(length));
  if (safeLength === 0) return [];

  const sorted = [...monitorResults].sort(
    (a, b) => b.createdAt.getTime() - a.createdAt.getTime(),
  );
  const recent = sorted.slice(0, safeLength);

  return Array.from({ length: safeLength }, (_, i) => recent[safeLength - 1 - i]);
};

