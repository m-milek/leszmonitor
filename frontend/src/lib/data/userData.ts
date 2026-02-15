import { BACKEND_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";
import { TeamRole, type User } from "@/lib/types.ts";

export const fetchUser = async (username: string): Promise<User> => {
  const res = await authFetch(`${BACKEND_URL}/users/${username}`);

  if (!res.ok) {
    throw new Error("Failed to fetch user");
  }

  return (await res.json()) as User;
};

export const fetchAllUsers = async (): Promise<User[]> => {
  const res = await authFetch(`${BACKEND_URL}/users`);

  if (!res.ok) {
    throw new Error("Failed to fetch users");
  }

  return (await res.json()) as User[];
};

export interface AddUserToTeamPayload {
  username: string;
  role: TeamRole;
}

export const addUserToTeam = async (
  teamId: string,
  { username, role }: AddUserToTeamPayload,
): Promise<void> => {
  const res = await authFetch(`${BACKEND_URL}/teams/${teamId}/members`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, role }),
  });

  if (!res.ok) {
    throw new Error("Failed to add user to team");
  }
};

export interface RemoveUserFromTeamPayload {
  username: string;
}

export const removeUserFromTeam = async (
  teamId: string,
  payload: RemoveUserFromTeamPayload,
): Promise<void> => {
  const res = await authFetch(`${BACKEND_URL}/teams/${teamId}/members`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    throw new Error("Failed to remove user from team");
  }
};

export interface RegisterUserPayload {
  username: string;
  password: string;
}

export const registerUser = async (
  payload: RegisterUserPayload,
): Promise<void> => {
  const res = await fetch(`${BACKEND_URL}/auth/register`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    throw new Error("Failed to register user");
  }
};
