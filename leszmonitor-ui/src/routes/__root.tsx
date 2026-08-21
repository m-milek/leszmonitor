import {
  Outlet,
  createRootRouteWithContext,
  redirect,
} from "@tanstack/react-router";
import type { QueryClient } from "@tanstack/react-query";
import { Providers } from "@/components/leszmonitor/Providers.tsx";
import { GlobalNotFound } from "@/components/leszmonitor/GlobalNotFound.tsx";
import { isJwtValid } from "@/lib/jwt.ts";
import { getCookie } from "@/lib/cookies.ts";

export interface RouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: () => (
    <Providers>
      <Outlet />
    </Providers>
  ),
  beforeLoad: async ({ location }) => {
    if (
      location.pathname === "/login" ||
      location.pathname === "/login/" ||
      location.pathname === "/register" ||
      location.pathname === "/register/"
    ) {
      return;
    }

    // Get token from cookie
    const token = getCookie("LOGIN_TOKEN");

    // Redirect to login if token is missing or invalid
    if (!token || !isJwtValid(token)) {
      throw redirect({ to: "/login" });
    }
  },
  notFoundComponent: GlobalNotFound,
});
