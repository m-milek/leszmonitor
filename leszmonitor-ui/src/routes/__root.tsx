import {
  HeadContent,
  Outlet,
  createRootRouteWithContext,
  redirect,
} from "@tanstack/react-router";
import type { QueryClient } from "@tanstack/react-query";
import { Providers } from "@/app/providers/AppProviders.tsx";
import { GlobalNotFound } from "@/components/common/GlobalNotFound.tsx";
import { isJwtValid } from "@/lib/jwt.ts";
import { readTokenSync } from "@/features/auth/lib/token.ts";

export interface RouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  head: () => ({
    meta: [{ title: "Leszmonitor" }],
  }),
  component: () => (
    <Providers>
      <HeadContent />
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

    const token = readTokenSync();

    // Redirect to login if token is missing or invalid
    if (!token || !isJwtValid(token)) {
      throw redirect({ to: "/login" });
    }
  },
  notFoundComponent: GlobalNotFound,
});
