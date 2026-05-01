import { BACKEND_API_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";
import type { Monitor } from "@/lib/types.ts";

export const getMonitors = async (projectId: string): Promise<Monitor[]> => {
  const res = await authFetch(`${BACKEND_API_URL}/projects/${projectId}/monitors`);

  return res.json();
};

export const getMonitorBySlug = async (
  projectId: string,
  monitorSlug: string,
): Promise<Monitor> => {
  const res = await authFetch(
    `${BACKEND_API_URL}/projects/${projectId}/monitors/${monitorSlug}`,
  );

  return res.json();
};

export const createMonitor = async (monitorData: Omit<Monitor, "id">) => {
  const res = await authFetch(
    `${BACKEND_API_URL}/projects/${monitorData.projectId}/monitors`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(monitorData),
    },
  );

  return res.json();
};
