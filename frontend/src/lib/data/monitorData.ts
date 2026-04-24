import { BACKEND_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";
import type { Monitor } from "@/lib/types.ts";

export const getMonitors = async (projectId: string): Promise<Monitor[]> => {
  const res = await authFetch(`${BACKEND_URL}/projects/${projectId}/monitors`);

  if (!res.ok) {
    throw new Error("Failed to fetch monitors");
  }

  return res.json();
};

export const createMonitor = async (monitorData: Omit<Monitor, "id">) => {
  const res = await authFetch(
    `${BACKEND_URL}/projects/${monitorData.projectId}/monitors`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(monitorData),
    },
  );

  if (!res.ok) {
    throw new Error("Failed to create monitor");
  }

  return res.json();
};
