import type { LoginPayload, LoginResponse } from "@/lib/types.ts";
import { BACKEND_URL } from "@/lib/consts.ts";

export const fetchLoginToken = async ({
  username,
  password,
}: LoginPayload): Promise<LoginResponse> => {
  const res = await fetch(`${BACKEND_URL}/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  });

  if (!res.ok) {
    throw new Error("Login failed");
  }

  const data = (await res.json()) as LoginResponse;
  if (!data.jwt) {
    throw new Error("Invalid response from server");
  }

  return data;
};
