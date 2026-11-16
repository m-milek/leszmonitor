import { SERVER_API_URL } from "@/lib/consts";
import { authFetch } from "@/lib/api-client";
import { readToken } from "@/features/auth/lib/token";
import type { PublicMonitor, StatusPage } from "@/features/status-pages/types";

export interface StatusPagePayload {
  slug: string;
  title: string;
  description: string;
  groupMonitorsByTag: boolean;
  requireAuth: boolean;
  enabled: boolean;
  timeWindowSeconds: number;
  tagIds: string[];
}

const mapStatusPage = (page: StatusPage): StatusPage => {
  return {
    ...page,
    createdAt: new Date(page.createdAt),
    updatedAt: new Date(page.updatedAt),
  };
};

const getAll = async (): Promise<StatusPage[]> => {
  const res = await authFetch(`${SERVER_API_URL}/status-pages`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const pages = (await res.json()) as StatusPage[];
  return pages.map(mapStatusPage);
};

const create = async (payload: StatusPagePayload): Promise<StatusPage> => {
  const res = await authFetch(`${SERVER_API_URL}/status-pages`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  const page = (await res.json()) as StatusPage;
  return mapStatusPage(page);
};

const update = async (
  id: string,
  payload: StatusPagePayload,
): Promise<StatusPage> => {
  const res = await authFetch(`${SERVER_API_URL}/status-pages/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  const page = (await res.json()) as StatusPage;
  return mapStatusPage(page);
};

const remove = async (id: string): Promise<void> => {
  await authFetch(`${SERVER_API_URL}/status-pages/${id}`, {
    method: "DELETE",
  });
};

const exportYaml = async (id: string): Promise<string> => {
  const res = await authFetch(`${SERVER_API_URL}/status-pages/${id}/export`, {
    method: "GET",
  });

  return res.text();
};

const publicFetch = async (url: string): Promise<Response> => {
  const token = await readToken();
  const res = await fetch(url, {
    headers: token ? { Authorization: `Bearer ${token}` } : undefined,
  });

  if (!res.ok) {
    throw new Error(`Failed to fetch ${url}: ${res.status}`, {
      cause: res.status,
    });
  }

  return res;
};

const getPublic = async (slug: string): Promise<StatusPage> => {
  const res = await publicFetch(
    `${SERVER_API_URL}/public/status-pages/${slug}`,
  );
  const page = (await res.json()) as StatusPage;
  return mapStatusPage(page);
};

const getPublicMonitors = async (slug: string): Promise<PublicMonitor[]> => {
  const res = await publicFetch(
    `${SERVER_API_URL}/public/status-pages/${slug}/monitors`,
  );
  const monitors = (await res.json()) as PublicMonitor[];
  return monitors.map((monitor) => ({
    ...monitor,
    tags: monitor.tags.map((tag) => ({
      ...tag,
      createdAt: new Date(tag.createdAt),
      updatedAt: new Date(tag.updatedAt),
    })),
    results: monitor.results.map((result) => ({
      ...result,
      createdAt: new Date(result.createdAt),
    })),
  }));
};

export const StatusPagesApi = {
  getAll,
  create,
  update,
  remove,
  exportYaml,
  public: {
    get: getPublic,
    getMonitors: getPublicMonitors,
  },
};
