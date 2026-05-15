import { authFetch } from "@/lib/data/utils.ts";
import { BACKEND_API_URL } from "@/lib/consts.ts";
import type { MonitorResult, Pagination } from "@/lib/types.ts";

export const getLatestMonitorResultByMonitorId = async (
  monitorId: string,
): Promise<MonitorResult | null> => {
  const res = await authFetch(
    `${BACKEND_API_URL}/monitors/${monitorId}/results/latest`,
  );

  if (res.status === 404) return null;
  if (!res.ok)
    throw new Error(`Failed to fetch latest result for ${monitorId}`);

  return res.json();
};

export const getMonitorResultsByMonitorId = async (
  monitorId: string,
  pagination: Pagination,
): Promise<MonitorResult[] | null> => {
  console.log(
    `Fetching results for monitor ${monitorId} with pagination:`,
    pagination,
  );
  const queryParams = new URLSearchParams({
    page: pagination.page.toString(),
    per_page: pagination.perPage.toString(),
  });
  const res = await authFetch(
    `${BACKEND_API_URL}/monitors/${monitorId}/results?${queryParams.toString()}`,
  );

  if (res.status === 404) return null;
  if (!res.ok) throw new Error(`Failed to fetch results for ${monitorId}`);

  return res.json();
};
