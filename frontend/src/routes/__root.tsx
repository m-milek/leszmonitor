import { Outlet, createRootRoute, redirect } from "@tanstack/react-router";
import { Providers } from "@/components/leszmonitor/Providers.tsx";

export const Route = createRootRoute({
  component: () => (
    <Providers>
      <Outlet />
    </Providers>
  ),
  beforeLoad: async ({ location }) => {
    if (location.pathname === "/login" || location.pathname === "/login/") {
      return;
    }

    const token = await cookieStore.get("LOGIN_TOKEN");
    if (!token) {
      throw redirect({ to: "/login" });
    }
  },
});
