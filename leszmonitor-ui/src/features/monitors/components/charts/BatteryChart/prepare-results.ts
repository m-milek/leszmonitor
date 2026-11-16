import type { MonitorStatus } from "@/features/monitors/types";

export interface BatteryChartResult {
  status: MonitorStatus;
  createdAt: Date;
}

export const prepareResults = <T extends BatteryChartResult>(
  monitorResults: T[],
  length: number,
): Array<T | undefined> => {
  const safeLength = Math.max(0, Math.floor(length));
  if (safeLength === 0) return [];

  const sorted = [...monitorResults].sort(
    (a, b) => b.createdAt.getTime() - a.createdAt.getTime(),
  );
  const recent = sorted.slice(0, safeLength);

  return Array.from(
    { length: safeLength },
    (_, i) => recent[safeLength - 1 - i],
  );
};
