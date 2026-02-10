import { getLoginToken } from "@/lib/utils.ts";
import type { Group } from "@/lib/types.ts";
import { BACKEND_URL } from "@/lib/consts.ts";

export async function fetchGroups(teamName: string): Promise<Group[] | null> {
  const token = await getLoginToken();
  if (!token) {
    return null;
  }

  const res = await fetch(`${BACKEND_URL}/teams/${teamName}/groups`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) {
    throw new Error("Failed to fetch groups");
  }

  return await res.json();
}
