import { UsersApi } from "@/features/users/users-api";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarRail,
} from "@/components/ui/sidebar";

import {
  LucideBookText,
  LucideHome,
  LucideLogs,
  LucideSearch,
  LucideSettings,
  LucideTag,
} from "lucide-react";
import { useAppStore } from "@/app/store";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useRef } from "react";
import { jwtDecode } from "jwt-decode";
import type { JwtClaims } from "@/lib/jwt";
import { AppSidebarFooter } from "@/components/layout/AppSidebarFooter";
import { readToken } from "@/features/auth/lib/token";
import { AppSidebarHeader } from "@/components/layout/AppSidebarHeader";
import { SidebarButton } from "@/components/layout/SidebarButton";

export const AppSidebar = () => {
  const { username, setUsername, user, setUser } = useAppStore();
  const hasInitialized = useRef(false);

  useEffect(() => {
    if (hasInitialized.current) return;
    hasInitialized.current = true;

    const getTokenAndExtractUsername = async () => {
      const token = await readToken();
      if (token) {
        const claims = jwtDecode(token) as JwtClaims;
        if (claims?.username) {
          setUsername(claims.username);
        }
      }
    };

    getTokenAndExtractUsername();
  }, [setUsername]);

  const { data: userData } = useQuery({
    queryKey: ["user", username],
    queryFn: () => UsersApi.get(username!),
    enabled: !!username,
    staleTime: 5 * 60 * 1000,
  });

  useEffect(() => {
    if (userData) {
      setUser(userData);
    }
  }, [userData, setUser]);

  return (
    <Sidebar variant="inset">
      <AppSidebarHeader />

      <SidebarContent className="p-2">
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu className="gap-1">
              <SidebarButton
                icon={<LucideHome />}
                href="/monitors"
                label="Home"
              />
              <SidebarButton icon={<LucideTag />} href="/tags" label="Tags" />
              {(user?.role === "owner" || user?.role === "admin") && (
                <SidebarButton
                  icon={<LucideLogs />}
                  href="/audit-log"
                  label="Audit Log"
                />
              )}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupLabel>Administration</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu className="gap-1">
              <SidebarButton
                icon={<LucideSettings />}
                href={`/admin`}
                label="Administration"
              />
              <SidebarButton
                icon={<LucideSettings />}
                href={`/user/${user?.username}/settings`}
                label="Settings"
              />
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupLabel>Help</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu className="gap-1">
              <SidebarButton
                icon={<LucideSearch />}
                href="/search"
                label="Search"
              />
              <SidebarButton
                icon={<LucideBookText />}
                href="/docs"
                label="Documentation"
              />
              <SidebarButton
                icon={<LucideLogs />}
                href="/_logdy"
                label="Logs"
              />
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter>{user && <AppSidebarFooter user={user} />}</SidebarFooter>

      <SidebarRail />
    </Sidebar>
  );
};
