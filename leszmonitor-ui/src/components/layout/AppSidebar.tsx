import { UsersApi } from "@/features/users/users-api";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  useSidebar
} from "@/components/ui/sidebar";

import {
  LucideActivity,
  LucideDiagramFlowVertical,
  LucideLogs,
  LucidePanelsTopLeft,
  LucidePlusCircle,
  LucideSearch,
  LucideSettings,
  LucideTag
} from "lucide-react";
import { useAppStore } from "@/app/store";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useRef } from "react";
import { jwtDecode } from "jwt-decode";
import type { JwtClaims } from "@/lib/jwt";
import { AppSidebarFooter, AppSidebarFooterSkeleton } from "@/components/layout/AppSidebarFooter";
import { readToken } from "@/features/auth/lib/token";
import { AppSidebarHeader } from "@/components/layout/AppSidebarHeader";
import { SidebarButton } from "@/components/layout/SidebarButton";
import { Link } from "@tanstack/react-router";

export const AppSidebar = () => {
  const { username, setUsername, user, setUser } = useAppStore();
  const { setOpenMobile } = useSidebar();
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
    <Sidebar variant="inset" collapsible="icon">
      <AppSidebarHeader />

      <SidebarContent className="p-2 group-data-[collapsible=icon]:p-0">
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu className="gap-1">
              <SidebarMenuItem className="mb-2">
                <SidebarMenuButton
                  render={<Link to="/monitors/new" draggable={false} />}
                  tooltip="New Monitor"
                  onClick={() => setOpenMobile(false)}
                  className="bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground"
                >
                  <LucidePlusCircle />
                  <span>New Monitor</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarButton
                icon={<LucideActivity />}
                href="/monitors"
                label="Monitors"
              />
              <SidebarButton icon={<LucideTag />} href="/tags" label="Tags" />
              <SidebarButton
                icon={<LucidePanelsTopLeft />}
                href="/statusPages"
                label="Status Pages"
              />
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        {(user?.role === "owner" || user?.role === "admin") && (
          <SidebarGroup>
            <SidebarGroupLabel>Administration</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu className="gap-1">
                <SidebarButton
                  icon={<LucideSettings />}
                  href="/admin"
                  label="Instance Settings"
                />
                <SidebarButton
                  icon={<LucideLogs />}
                  href="/audit-log"
                  label="Audit Log"
                />
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        )}
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
                icon={<LucideDiagramFlowVertical />}
                href="/docs"
                label="API Reference"
              />
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter>
        {user ? <AppSidebarFooter user={user} /> : <AppSidebarFooterSkeleton />}
      </SidebarFooter>

      <SidebarRail />
    </Sidebar>
  );
};
