import { createFileRoute, redirect } from "@tanstack/react-router";
import { jwtDecode } from "jwt-decode";
import type { JwtClaims } from "@/lib/types.ts";

export const Route = createFileRoute("/_authenticated/")({
  beforeLoad: async () => {
    const token = await cookieStore.get("LOGIN_TOKEN");

    if (token?.value) {
      const claims = jwtDecode(token.value) as JwtClaims;

      if (claims.username) {
        throw redirect({
          to: "/$teamId",
          params: { teamId: claims.username },
        });
      }
    }
  },
  component: () => null,
});
