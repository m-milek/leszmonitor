import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar.tsx";

import {
  LayoutDashboardIcon,
  LucideActivity,
  LucideBookText,
  LucideFolderOpen,
  LucidePlusCircle,
  LucideSearch,
  LucideSettings,
  LucideUsers,
} from "lucide-react";
import { Logo } from "@/components/leszmonitor/Logo.tsx";
import { useAtom, useAtomValue } from "jotai";
import { teamAtom, userAtom, usernameAtom } from "@/lib/atoms.ts";
import { useQuery } from "@tanstack/react-query";
import { fetchUser } from "@/lib/fetchUser.ts";
import { useEffect, useRef } from "react";
import { jwtDecode } from "jwt-decode";
import type { JwtClaims } from "@/lib/types.ts";
import { TeamSelector } from "@/components/leszmonitor/sidebar/TeamSelector.tsx";
import { Link } from "@tanstack/react-router";
import { AppSidebarFooter } from "@/components/leszmonitor/sidebar/AppSidebarFooter.tsx";

interface SidebarButtonProps {
  icon: React.ReactNode;
  href: string;
  label: string;
  variant?: "default" | "primary";
}

const SidebarButton = ({
  icon,
  href,
  label,
  variant = "default",
}: SidebarButtonProps) => {
  const onClick = () => {
    console.log("navigating to", href);
  };

  const className =
    variant === "primary"
      ? "bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground"
      : "";

  return (
    <SidebarMenuItem>
      <SidebarMenuButton onClick={onClick} className={className}>
        <Link to={href} className="flex items-center space-x-2">
          {icon}
          <span>{label}</span>
        </Link>
      </SidebarMenuButton>
    </SidebarMenuItem>
  );
};

export const AppSidebar = () => {
  const [username, setUsername] = useAtom(usernameAtom);
  const [user, setUser] = useAtom(userAtom);
  const hasInitialized = useRef(false);

  const team = useAtomValue(teamAtom);

  useEffect(() => {
    if (hasInitialized.current) return;
    hasInitialized.current = true;

    const getTokenAndExtractUsername = async () => {
      const token = await cookieStore.get("LOGIN_TOKEN");
      if (token?.value) {
        const claims = jwtDecode(token.value) as JwtClaims;
        if (claims?.username) {
          setUsername(claims.username);
        }
      }
    };

    getTokenAndExtractUsername();
  }, [setUsername]);

  const { data: userData } = useQuery({
    queryKey: ["user", username],
    queryFn: () => fetchUser(username!),
    enabled: !!username,
    staleTime: 5 * 60 * 1000,
  });

  useEffect(() => {
    if (userData) {
      setUser(userData);
    }
  }, [userData, setUser]);

  if (!username || !team || !user) {
    return null;
  }

  return (
    <Sidebar variant="inset">
      <SidebarHeader>
        <div className="p-2">
          <Logo />
        </div>
        <SidebarMenu>
          <SidebarMenuItem>
            <TeamSelector />
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarButton
                icon={<LucidePlusCircle />}
                href="/new-monitor"
                label="New Monitor"
                variant="primary"
              />
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupLabel>This Workspace</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarButton
                icon={<LayoutDashboardIcon />}
                href={`/${team.displayId}/dashboard`}
                label="Dashboard"
              />
              <SidebarButton
                icon={<LucideActivity />}
                href={`/${team.displayId}/monitors`}
                label="Monitors"
              />
              <SidebarButton
                icon={<LucideFolderOpen />}
                href={`/${team.displayId}/groups`}
                label="Groups"
              />
              <SidebarButton
                icon={<LucideUsers />}
                href={`/${team.displayId}/team`}
                label="Team"
              />
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupLabel>Help</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarButton
                icon={<LucideBookText />}
                href="/docs"
                label="Documentation"
              />
              <SidebarButton
                icon={<LucideSearch />}
                href="/search"
                label="Search"
              />
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarButton
              icon={<LucideSettings />}
              href="/settings"
              label="Settings"
            />
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <AppSidebarFooter user={user} />
      </SidebarFooter>
    </Sidebar>
  );
};
