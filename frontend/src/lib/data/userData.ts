import { BACKEND_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";
import { OrgRole, type User } from "@/lib/types.ts";

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

export interface AddUserToOrgPayload {
  username: string;
  role: OrgRole;
}

export const addUserToOrg = async (
  orgId: string,
  { username, role }: AddUserToOrgPayload,
): Promise<void> => {
  const res = await authFetch(`${BACKEND_URL}/orgs/${orgId}/members`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, role }),
  });

  if (!res.ok) {
    throw new Error("Failed to add user to org");
  }
};

export interface RemoveUserFromOrgPayload {
  username: string;
}

export const removeUserFromOrg = async (
  orgId: string,
  payload: RemoveUserFromOrgPayload,
): Promise<void> => {
  const res = await authFetch(`${BACKEND_URL}/orgs/${orgId}/members`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    throw new Error("Failed to remove user from org");
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
