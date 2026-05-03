import { authFetch } from "@/lib/data/utils.ts";
import { BACKEND_API_URL } from "@/lib/consts.ts";
import type { MonitorResult } from "@/lib/types.ts";

export const getLatestMonitorResultById = async (
  monitorId: string,
): Promise<MonitorResult> => {
  const res = await authFetch(
    `${BACKEND_API_URL}/monitors/${monitorId}/results/latest`,
  );

  return res.json();
};
