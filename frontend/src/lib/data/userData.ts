import { BACKEND_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";
import type { User } from "@/lib/types.ts";

export const fetchUser = async (username: string): Promise<User> => {
  const res = await authFetch(`${BACKEND_URL}/users/${username}`);

  if (!res.ok) {
    throw new Error("Failed to fetch user");
  }

  return (await res.json()) as User;
};
