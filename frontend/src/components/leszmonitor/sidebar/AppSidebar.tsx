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
import { LeszmonitorLogo } from "@/components/leszmonitor/ui/LeszmonitorLogo.tsx";
import { useAtom, useAtomValue } from "jotai";
import { projectAtom, userAtom, usernameAtom } from "@/lib/atoms.ts";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useRef } from "react";
import { jwtDecode } from "jwt-decode";
import type { JwtClaims } from "@/lib/types.ts";
import { ProjectSelector } from "@/components/leszmonitor/sidebar/OrgSelector.tsx";
import { Link } from "@tanstack/react-router";
import { AppSidebarFooter } from "@/components/leszmonitor/sidebar/AppSidebarFooter.tsx";
import { fetchUser } from "@/lib/data/userData.ts";
import { getProjects } from "@/lib/data/projectData.ts";
import { Skeleton } from "@/components/ui/skeleton.tsx";

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
  const [username, setUsername] = useAtom(usernameAtom);
  const [user, setUser] = useAtom(userAtom);
  const hasInitialized = useRef(false);

  const project = useAtomValue(projectAtom);

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
    queryFn: getProjects,
  });

  useEffect(() => {
    if (userData) {
      setUser(userData);
    }
  }, [userData, setUser]);

  return (
    <Sidebar variant="inset">
      <SidebarHeader>
        <div className="p-2">
          <Link to={"/"}>
            <LeszmonitorLogo />
          </Link>
        </div>
        <SidebarMenu>
          <SidebarMenuItem>
            {projectsData ? (
              <ProjectSelector projects={projectsData} />
            ) : (
              <Skeleton className="h-8 w-full" />
            )}
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent>
        {project && (
          <>
            <SidebarGroup>
              <SidebarGroupContent>
                <SidebarMenu>
                  <SidebarButton
                    icon={<LucidePlusCircle />}
                    href={`/projects/${project.slug}/monitors/new`}
                    label="New Monitor"
                    variant="primary"
                  />
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
            <SidebarGroup>
              <SidebarGroupLabel>This Project</SidebarGroupLabel>
              <SidebarGroupContent>
                <SidebarMenu>
                  <SidebarButton
                    icon={<LayoutDashboardIcon />}
                    href={`/projects/${project.slug}/dashboard`}
                    label="Dashboard"
                  />
                  <SidebarButton
                    icon={<LucideActivity />}
                    href={`/projects/${project.slug}/monitors`}
                    label="Monitors"
                  />
                  <SidebarButton
                    icon={<LucideUsers />}
                    href={`/projects/${project.slug}/members`}
                    label="Members"
                  />
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          </>
        )}
        <SidebarGroup>
          <SidebarGroupLabel>Workspaces</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarButton
                icon={<LucideFolderOpen />}
                href="/projects"
                label="All Projects"
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
              <SidebarButton
                icon={<LucideSettings />}
                href={`/user/${user?.username}/settings`}
                label="Settings"
              />
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter>{user && <AppSidebarFooter user={user} />}</SidebarFooter>
    </Sidebar>
  );
};
