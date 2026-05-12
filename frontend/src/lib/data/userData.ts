import { BACKEND_API_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";
import type { ApiError, User } from "@/lib/types.ts";

export const fetchUser = async (username: string): Promise<User> => {
  const res = await authFetch(`${BACKEND_API_URL}/users/${username}`);

  return (await res.json()) as User;
};

export const fetchAllUsers = async (): Promise<User[]> => {
  const res = await authFetch(`${BACKEND_API_URL}/users`);

  return (await res.json()) as User[];
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
