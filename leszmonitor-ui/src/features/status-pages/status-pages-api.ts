import { SERVER_API_URL } from "@/lib/consts";
import { authFetch } from "@/lib/api-client";
import type { StatusPage } from "@/features/status-pages/types";

export interface StatusPagePayload {
  slug: string;
  title: string;
  description: string;
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

export const StatusPagesApi = {
  getAll,
  create,
};
