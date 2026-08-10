import { BACKEND_API_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";
import type { ApiError, User } from "@/lib/types.ts";

export const getUser = async (username: string): Promise<User> => {
  const res = await authFetch(`${BACKEND_API_URL}/users/${username}`);

  const user = (await res.json()) as User;
  return mapUser(user);
};

export const getAllUsers = async (): Promise<User[]> => {
  const res = await authFetch(`${BACKEND_API_URL}/users`);

  const users = (await res.json()) as User[];
  return users.map(mapUser);
};

export interface RegisterUserPayload {
  username: string;
  password: string;
}

export const registerUser = async (
  payload: RegisterUserPayload,
): Promise<void> => {
  const response = await fetch(`${BACKEND_API_URL}/auth/register`, {
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

const mapUser = (user: User): User => {
  return {
    ...user,
    createdAt: new Date(user.createdAt),
    updatedAt: new Date(user.updatedAt),
  };
};

export const removeUser = async (username: string): Promise<void> => {
  // Mock user removal data layer function
  console.log(`Mock removing user: ${username}`);
  return Promise.resolve();
};
