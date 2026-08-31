import { authFetch } from "@/lib/data/utils.ts";
import { SERVER_API_URL } from "@/lib/consts.ts";

export interface AverageLatencyParams {
  from: Date;
  to?: Date;
}

export interface AverageLatencyResponse {
  averageLatency: number;
}

export const getAverageLatencyByMonitorId = async (
  monitorId: string,
  { from, to = new Date(Date.now()) }: AverageLatencyParams,
): Promise<AverageLatencyResponse> => {
  const queryParams = new URLSearchParams({
    from: from.toISOString(),
    to: to.toISOString(),
  });
  const res = await authFetch(
    `${SERVER_API_URL}/monitors/${monitorId}/stats/average-latency?${queryParams.toString()}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    },
  );

  if (!res.ok)
    throw new Error(`Failed to fetch average latency for ${monitorId}`);

  return res.json();
};
