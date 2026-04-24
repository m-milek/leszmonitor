import { BACKEND_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";
import type { User } from "@/lib/types.ts";

export const fetchUser = async (username: string): Promise<User> => {
  const res = await authFetch(`${BACKEND_URL}/users/${username}`);

  return (await res.json()) as User;
};

export const fetchAllUsers = async (): Promise<User[]> => {
  const res = await authFetch(`${BACKEND_URL}/users`);

  return (await res.json()) as User[];
};

export interface RegisterUserPayload {
  username: string;
  password: string;
}

export const registerUser = async (
  payload: RegisterUserPayload,
): Promise<void> => {
  await fetch(`${BACKEND_URL}/auth/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
};
