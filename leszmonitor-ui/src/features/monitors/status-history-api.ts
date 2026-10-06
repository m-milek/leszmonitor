import { authFetch } from "@/lib/api-client";
import { SERVER_API_URL } from "@/lib/consts";
import type { MonitorStatusPeriod } from "@/features/monitors/types";
import type { Pagination } from "@/lib/types";

const getPage = async (
  monitorId: string,
  pagination: Pagination,
  from: Date,
): Promise<MonitorStatusPeriod[]> => {
  const queryParams = new URLSearchParams({
    page: pagination.page.toString(),
    per_page: pagination.perPage.toString(),
    from: from.toISOString(),
    to: new Date().toISOString(),
  });
  const res = await authFetch(
    `${SERVER_API_URL}/monitors/${monitorId}/statusHistory?${queryParams.toString()}`,
  );

  if (!res.ok)
    throw new Error(`Failed to fetch status history for ${monitorId}`);

  return (await res.json()).map(mapMonitorStatusPeriod);
};

const mapMonitorStatusPeriod = (
  period: MonitorStatusPeriod,
): MonitorStatusPeriod => {
  return {
    ...period,
    startedAt: new Date(period.startedAt),
    endedAt: period.endedAt ? new Date(period.endedAt) : null,
  };
};

export const statusHistoryApi = {
  getPage,
};
