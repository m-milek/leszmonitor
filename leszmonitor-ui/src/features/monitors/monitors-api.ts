import { resultsApi } from "@/features/monitors/results-api";
import { statsApi } from "@/features/monitors/stats-api";
import { statusHistoryApi } from "@/features/monitors/status-history-api";
import { SERVER_API_URL } from "@/lib/consts";
import { authFetch } from "@/lib/api-client";
import type {
  Monitor,
  MonitorCreatePayload,
  MonitorUpdatePayload,
} from "@/features/monitors/types";

const getAll = async (): Promise<Monitor[]> => {
  const res = await authFetch(`${SERVER_API_URL}/monitors`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return (await res.json()) as Monitor[];
};

const getBySlug = async (monitorSlug: string): Promise<Monitor> => {
  const res = await authFetch(`${SERVER_API_URL}/monitors/${monitorSlug}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return (await res.json()) as Monitor;
};

const create = async (monitorData: MonitorCreatePayload): Promise<Monitor> => {
  const res = await authFetch(`${SERVER_API_URL}/monitors`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(monitorData),
  });

  return (await res.json()) as Monitor;
};

const update = async (
  monitorId: string,
  monitorData: MonitorUpdatePayload,
): Promise<Monitor> => {
  const res = await authFetch(`${SERVER_API_URL}/monitors/${monitorId}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(monitorData),
  });

  return (await res.json()) as Monitor;
};

const updateState = async (
  monitorId: string,
  newState: string,
): Promise<Monitor> => {
  const res = await authFetch(`${SERVER_API_URL}/monitors/${monitorId}/state`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ newState }),
  });

  return (await res.json()) as Monitor;
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

const exportYaml = async (monitorId: string): Promise<string> => {
  const res = await authFetch(
    `${SERVER_API_URL}/monitors/${monitorId}/export`,
    {
      method: "GET",
    },
  );

  return res.text();
};

export const MonitorsApi = {
  getAll,
  getBySlug,
  create,
  update,
  updateState,
  remove,
  run,
  exportYaml,
  results: resultsApi,
  stats: statsApi,
  statusHistory: statusHistoryApi,
};
