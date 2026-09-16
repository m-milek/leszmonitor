import { readToken } from "@/features/auth/lib/token";

export const authFetch = async (url: string, options?: RequestInit) => {
  const token = await readToken();
  if (!token) {
    throw new Error("No login token found");
  }

  const res = await fetch(url, {
    ...options,
    headers: {
      ...options?.headers,
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) {
    throw new Error(`Failed to fetch ${url}: ${res.statusText}`);
  }

  return res;
};
