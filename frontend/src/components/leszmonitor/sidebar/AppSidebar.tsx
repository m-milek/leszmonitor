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
} from "@/components/ui/sidebar.tsx";

import {
  LucideBookText,
  LucideHome,
  LucideSearch,
  LucideSettings,
  LucideUsers,
} from "lucide-react";
import { useAppStore } from "@/lib/store.ts";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useRef } from "react";
import { jwtDecode } from "jwt-decode";
import type { JwtClaims } from "@/lib/types.ts";
import { Link } from "@tanstack/react-router";
import { AppSidebarFooter } from "@/components/leszmonitor/sidebar/AppSidebarFooter.tsx";
import { fetchUser } from "@/lib/data/userData.ts";
import { getProjects } from "@/lib/data/projectData.ts";
import { AppSidebarHeader } from "@/components/leszmonitor/sidebar/AppSidebarHeader.tsx";
import { ProjectMenu } from "@/components/leszmonitor/sidebar/ProjectMenu.tsx";

interface SidebarButtonProps {
  icon: React.ReactNode;
  href: string;
  label: string;
  variant?: "default" | "primary";
}

export const SidebarButton = ({
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
      ? "bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground active:bg-primary/90 active:text-primary-foreground cursor-pointer"
      : "cursor-pointer";

  return (
    <SidebarMenuItem>
      <Link to={href} className="flex items-center space-x-2 cursor-pointer">
        <SidebarMenuButton onClick={onClick} className={className}>
          {icon}
          <span>{label}</span>
        </SidebarMenuButton>
      </Link>
    </SidebarMenuItem>
  );
};

export const AppSidebar = () => {
  const { username, setUsername, user, setUser } = useAppStore();
  const hasInitialized = useRef(false);

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

  const { data: projectsData } = useQuery({
    queryKey: ["projects"],
    queryFn: () => getProjects(),
  });

  useEffect(() => {
    if (userData) {
      setUser(userData);
    }
  }, [userData, setUser]);

  return (
    <Sidebar variant="inset">
      <AppSidebarHeader />

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarButton
                icon={<LucideHome />}
                href="/projects"
                label="Home"
              />
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupContent>
            <ProjectMenu projects={projectsData ?? []} />
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupLabel>Administration</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarButton
                icon={<LucideUsers />}
                href={`/users`}
                label="Users"
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
            <SidebarMenu>
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
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter>{user && <AppSidebarFooter user={user} />}</SidebarFooter>
    </Sidebar>
  );
};
