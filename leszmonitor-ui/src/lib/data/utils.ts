export const authFetch = async (url: string, options?: RequestInit) => {
  const token = await cookieStore.get("LOGIN_TOKEN");
  if (!token) {
    throw new Error("No login token found");
  }

  const res = await fetch(url, {
    ...options,
    headers: {
      ...options?.headers,
      Authorization: `Bearer ${token.value}`,
    },
  });

  if (!res.ok) {
    throw new Error(`Failed to fetch ${url}: ${res.statusText}`);
  }

  return res;
};
