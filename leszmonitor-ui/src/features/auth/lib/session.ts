import { jwtDecode } from "jwt-decode";
import { fetchLoginToken } from "@/features/auth/api/auth.ts";
import { getUser } from "@/features/users/api/users.ts";
import { isJwtClaims } from "@/lib/jwt.ts";
import { storeToken } from "@/features/auth/lib/token.ts";
import type { LoginPayload } from "@/features/auth/model/types.ts";
import type { User } from "@/features/users/model/types.ts";

export interface SessionSetters {
  setUsername: (username: string | null) => void;
  setUser: (user: User | null) => void;
}

// Exchanges credentials for a token, persists it and hydrates the app store.
// Returns false when the token carries unusable claims, so the caller can skip
// navigating without surfacing an error.
export const establishSession = async (
  credentials: LoginPayload,
  { setUsername, setUser }: SessionSetters,
): Promise<boolean> => {
  const loginResponse = await fetchLoginToken(credentials);

  storeToken(loginResponse.jwt);

  const claims = jwtDecode(loginResponse.jwt);
  if (!isJwtClaims(claims)) {
    console.error("Invalid JWT claims");
    return false;
  }

  setUsername(claims.username);

  const user = await getUser(claims.username);
  setUser(user);

  return true;
};
