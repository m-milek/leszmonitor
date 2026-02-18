import { authFetch } from "@/lib/data/utils.ts";
import { BACKEND_URL } from "@/lib/consts.ts";
import type { Monitor } from "@/lib/types.ts";

export const getMonitors = async (teamId: string): Promise<Monitor> => {
  const res = await authFetch(`${BACKEND_URL}/teams/${teamId}/monitors`);

  if (!res.ok) {
    throw new Error("Failed to fetch monitors");
  }

  return res.json();
};
