export interface JwtClaims {
  iat: number;
  exp: number;
  username: string;
}

export const isJwtClaims = (data: unknown): data is JwtClaims => {
  if (typeof data !== "object" || data === null) {
    return false;
  }
  const obj = data as Record<string, unknown>;
  return (
    typeof obj.username === "string" &&
    typeof obj.exp === "number" &&
    typeof obj.iat === "number"
  );
}

export const isJwtValid = (token: string): JwtClaims | null => {
  try {
    const tokenDecoded = atob(token.split(".")[1]);
    const tokenData = JSON.parse(tokenDecoded) as JwtClaims;

    if (!isJwtClaims(tokenData)) {
      return null;
    }

    if (tokenData.exp * 1000 < Date.now()) {
      return null;
    }

    return tokenData;
  } catch (error) {
    console.error("Invalid JWT:", error);
    return null;
  }
}