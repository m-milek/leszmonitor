import { BACKEND_URL } from "@/lib/consts.ts";
import type { Group, Timestamps } from "@/lib/types.ts";
import { authFetch } from "@/lib/data/utils.ts";

export type GroupInput = Omit<Group, "id" | keyof Timestamps>;

export const getGroups = async (teamId: string): Promise<Group[]> => {
  const res = await authFetch(`${BACKEND_URL}/teams/${teamId}/groups`);

  if (!res.ok) {
    throw new Error("Failed to fetch groups");
  }

  const groups = (await res.json()) as Group[];
  return groups.map((group) => ({
    ...group,
    createdAt: new Date(group.createdAt),
    updatedAt: new Date(group.updatedAt),
  }));
};

export const addGroup = async (
  teamId: string,
  group: GroupInput,
): Promise<Group> => {
  const res = await authFetch(`${BACKEND_URL}/teams/${teamId}/groups`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(group),
  });

  if (!res.ok) {
    throw new Error("Failed to add group");
  }

  const newGroup = (await res.json()) as Group;
  return {
    ...newGroup,
    createdAt: new Date(newGroup.createdAt),
    updatedAt: new Date(newGroup.updatedAt),
  };
};
