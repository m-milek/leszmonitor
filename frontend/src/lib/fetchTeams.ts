import { getLoginToken } from "@/lib/utils.ts";
import type { Team } from "@/lib/types.ts";
import { BACKEND_URL } from "@/lib/consts.ts";

export async function fetchTeams(): Promise<Team[]> {
  const token = await getLoginToken();
  if (!token) {
    return [];
  }

  const res = await fetch(`${BACKEND_URL}/teams`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) {
    throw new Error("Failed to fetch teams");
  }

  return await res.json();
}
