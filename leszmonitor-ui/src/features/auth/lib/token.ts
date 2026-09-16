import { getCookie, setCookie } from "@/lib/cookies";

export const LOGIN_TOKEN_COOKIE = "LOGIN_TOKEN";

const TOKEN_MAX_AGE_SECONDS = 24 * 60 * 60;

// Two readers on purpose: the router guard runs synchronously and cannot await,
// while everything else already lives in async code and uses the richer
// cookieStore API. Do not collapse them into one.
export const readTokenSync = (): string | null => getCookie(LOGIN_TOKEN_COOKIE);

export const readToken = async (): Promise<string | null> => {
  const token = await cookieStore.get(LOGIN_TOKEN_COOKIE);
  return token?.value || null;
};

export const storeToken = (jwt: string): void => {
  setCookie(LOGIN_TOKEN_COOKIE, jwt, {
    maxAge: TOKEN_MAX_AGE_SECONDS,
    path: "/",
    sameSite: "Lax",
  });
};

export const clearToken = async (): Promise<void> => {
  await cookieStore.delete(LOGIN_TOKEN_COOKIE);
};
