import { BACKEND_API_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";
import type {
  Monitor,
  MonitorCreatePayload,
  MonitorUpdatePayload,
} from "@/lib/types.ts";

const normalizeMonitor = (monitor: Monitor): Monitor => {
  if (typeof monitor.probeConfig === "string") {
    try {
      return {
        ...monitor,
        probeConfig: JSON.parse(monitor.probeConfig),
      } as Monitor;
    } catch {
      return monitor;
    }
  }

  return monitor;
};

export const getMonitorsByProjectSlug = async (
  projectSlug: string,
): Promise<Monitor[]> => {
  const query = new URLSearchParams({
    projectSlug,
  }).toString();
  const res = await authFetch(`${BACKEND_API_URL}/monitors?${query}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const monitors = (await res.json()) as Monitor[];
  return monitors.map(normalizeMonitor);
};

export const getMonitorBySlug = async (
  projectSlug: string,
  monitorSlug: string,
): Promise<Monitor> => {
  const res = await authFetch(
    `${BACKEND_API_URL}/projects/${projectSlug}/monitors/${monitorSlug}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    },
  );

  const monitor = (await res.json()) as Monitor;
  return normalizeMonitor(monitor);
};

export const createMonitor = async (monitorData: MonitorCreatePayload) => {
  const query = new URLSearchParams({
    projectSlug: monitorData.projectSlug,
  }).toString();
  const res = await authFetch(`${BACKEND_API_URL}/monitors?${query}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(monitorData),
  });

  return res.json();
};

export const updateMonitor = async (
  monitorId: string,
  monitorData: MonitorUpdatePayload,
) => {
  const res = await authFetch(
    `${BACKEND_API_URL}/projects/${monitorData.projectSlug}/monitors/${monitorId}`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(monitorData),
    },
  );

  return res.json();
};

export const updateMonitorState = async (
  monitorId: string,
  newState: string,
) => {
  const res = await authFetch(
    `${BACKEND_API_URL}/monitors/${monitorId}/state`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ newState }),
    },
  );

  return res.json();
};

export const deleteMonitor = async (monitorId: string) => {
  await authFetch(`${BACKEND_API_URL}/monitors/${monitorId}`, {
    method: "DELETE",
  });
};
