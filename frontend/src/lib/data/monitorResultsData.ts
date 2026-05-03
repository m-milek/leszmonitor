import { authFetch } from "@/lib/data/utils.ts";
import { BACKEND_API_URL } from "@/lib/consts.ts";
import type { MonitorResult } from "@/lib/types.ts";

export const getLatestMonitorResultById = async (
  monitorId: string,
): Promise<MonitorResult | null> => {
  const res = await authFetch(
    `${BACKEND_API_URL}/monitors/${monitorId}/results/latest`,
  );

  if (res.status === 404) return null;
  if (!res.ok) throw new Error(`Failed to fetch result for ${monitorId}`);

  return res.json();
};
