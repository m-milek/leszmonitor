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

  const team = (await res.json()) as Team;

  team.updatedAt = new Date(team.updatedAt);
  team.createdAt = new Date(team.createdAt);
  team.members = team.members.map((member) => ({
    ...member,
    createdAt: new Date(member.createdAt),
    updatedAt: new Date(member.updatedAt),
  }));

  return team;
}
