import { SERVER_API_URL } from "@/lib/consts";
import { authFetch } from "@/lib/api-client";
import type { ApiError } from "@/lib/types";
import type { User, UserRole } from "@/features/users/types";

const get = async (username: string): Promise<User> => {
  const res = await authFetch(`${SERVER_API_URL}/users/${username}`);

  const user = (await res.json()) as User;
  return mapUser(user);
};

const getAll = async (): Promise<User[]> => {
  const res = await authFetch(`${SERVER_API_URL}/users`);

  const users = (await res.json()) as User[];
  return users.map(mapUser);
};

export interface RegisterUserPayload {
  username: string;
  password: string;
}

const register = async (payload: RegisterUserPayload): Promise<void> => {
  const response = await fetch(`${SERVER_API_URL}/auth/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    const errorData = (await response.json()) as ApiError;
    console.error(errorData);
    throw new Error("Failed to register user: " + errorData.error.message);
  }
};

const updateRole = async (username: string, role: UserRole): Promise<User> => {
  const res = await authFetch(`${SERVER_API_URL}/users/${username}/role`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ role }),
  });

  const user = (await res.json()) as User;
  return mapUser(user);
};

const mapUser = (user: User): User => {
  return {
    ...user,
    createdAt: new Date(user.createdAt),
    updatedAt: new Date(user.updatedAt),
  };
};

const remove = async (username: string): Promise<void> => {
  // Mock user removal data layer function
  console.log(`Mock removing user: ${username}`);
  return Promise.resolve();
};

export const UsersApi = {
  get,
  getAll,
  register,
  updateRole,
  remove,
};
