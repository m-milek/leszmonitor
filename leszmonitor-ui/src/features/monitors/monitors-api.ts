import { resultsApi } from "@/features/monitors/results-api";
import { statsApi } from "@/features/monitors/stats-api";
import { SERVER_API_URL } from "@/lib/consts";
import { authFetch } from "@/lib/api-client";
import type {
  Monitor,
  MonitorCreatePayload,
  MonitorUpdatePayload,
} from "@/features/monitors/types";

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

const getAll = async (): Promise<Monitor[]> => {
  const res = await authFetch(`${SERVER_API_URL}/monitors`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const monitors = (await res.json()) as Monitor[];
  return monitors.map(normalizeMonitor);
};

const getBySlug = async (monitorSlug: string): Promise<Monitor> => {
  const res = await authFetch(`${SERVER_API_URL}/monitors/${monitorSlug}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const monitor = (await res.json()) as Monitor;
  return normalizeMonitor(monitor);
};

const create = async (monitorData: MonitorCreatePayload) => {
  const res = await authFetch(`${SERVER_API_URL}/monitors`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(monitorData),
  });

  return res.json();
};

const update = async (monitorId: string, monitorData: MonitorUpdatePayload) => {
  const res = await authFetch(`${SERVER_API_URL}/monitors/${monitorId}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(monitorData),
  });

  return res.json();
};

const updateState = async (monitorId: string, newState: string) => {
  const res = await authFetch(`${SERVER_API_URL}/monitors/${monitorId}/state`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ newState }),
  });

  return res.json();
};

const remove = async (monitorId: string) => {
  await authFetch(`${SERVER_API_URL}/monitors/${monitorId}`, {
    method: "DELETE",
  });
};

const run = async (monitorId: string) => {
  const res = await authFetch(`${SERVER_API_URL}/monitors/${monitorId}/run`, {
    method: "POST",
  });

  return res.json();
};

export const MonitorsApi = {
  getAll,
  getBySlug,
  create,
  update,
  updateState,
  remove,
  run,
  results: resultsApi,
  stats: statsApi,
};
