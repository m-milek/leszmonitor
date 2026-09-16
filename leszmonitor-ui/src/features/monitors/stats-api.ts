import { authFetch } from "@/lib/api-client";
import { SERVER_API_URL } from "@/lib/consts";
import type { MonitorStats } from "@/features/monitors/types.ts";

export interface MonitorStatsParams {
  from: Date;
  to?: Date;
}

const get = async (
  monitorId: string,
  { from, to = new Date(Date.now()) }: MonitorStatsParams,
): Promise<MonitorStats> => {
  const queryParams = new URLSearchParams({
    from: from.toISOString(),
    to: to.toISOString(),
  });
  const res = await authFetch(
    `${SERVER_API_URL}/monitors/${monitorId}/stats?${queryParams.toString()}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    },
  );

  if (!res.ok)
    throw new Error(`Failed to fetch latency stats for ${monitorId}`);

  return res.json();
};

export const statsApi = {
  get,
};
