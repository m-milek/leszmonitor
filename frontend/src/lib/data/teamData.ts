import type { Team } from "@/lib/types.ts";
import { BACKEND_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";

export async function getTeam(teamName: string): Promise<Team | null> {
  const res = await authFetch(`${BACKEND_URL}/teams/${teamName}`);

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

export async function fetchTeams(): Promise<Team[]> {
  const res = await authFetch(`${BACKEND_URL}/teams`);

  if (!res.ok) {
    throw new Error("Failed to fetch teams");
  }

  return await res.json();
}
