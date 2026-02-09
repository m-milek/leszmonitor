import { getLoginToken } from "@/lib/utils.ts";
import type { Team } from "@/lib/types.ts";
import { BACKEND_URL } from "@/lib/consts.ts";

export async function fetchTeam(teamName: string): Promise<Team | null> {
  const token = await getLoginToken();
  if (!token) {
    return null;
  }

  const res = await fetch(`${BACKEND_URL}/teams/${teamName}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) {
    throw new Error("Failed to fetch team");
  }

  return await res.json();
}
