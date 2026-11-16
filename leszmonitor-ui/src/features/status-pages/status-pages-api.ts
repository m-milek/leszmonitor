import { SERVER_API_URL } from "@/lib/consts";
import { authFetch } from "@/lib/api-client";
import type { StatusPage } from "@/features/status-pages/types";

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

export const StatusPagesApi = {
  getAll,
};
