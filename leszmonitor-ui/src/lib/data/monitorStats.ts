import { authFetch } from "@/lib/data/utils.ts";
import { SERVER_API_URL } from "@/lib/consts.ts";

export interface LatencyStatsParams {
  from: Date;
  to?: Date;
}

export interface LatencyStatsResponse {
  averageLatency: number;
  minLatency: number;
  maxLatency: number;
}

export const getLatencyStatsByMonitorId = async (
  monitorId: string,
  { from, to = new Date(Date.now()) }: LatencyStatsParams,
): Promise<LatencyStatsResponse> => {
  const queryParams = new URLSearchParams({
    from: from.toISOString(),
    to: to.toISOString(),
  });
  const res = await authFetch(
    `${SERVER_API_URL}/monitors/${monitorId}/stats/latency?${queryParams.toString()}`,
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
