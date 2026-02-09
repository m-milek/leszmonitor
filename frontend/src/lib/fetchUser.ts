import { BACKEND_URL } from "@/lib/consts.ts";

export const fetchUser = async (username: string) => {
  const token = await cookieStore.get("LOGIN_TOKEN");
  if (!token) {
    return null;
  }

  const res = await fetch(`${BACKEND_URL}/users/${username}`, {
    headers: {
      Authorization: `Bearer ${token.value}`,
    },
  });

  if (!res.ok) {
    return null;
  }

  return await res.json();
};
