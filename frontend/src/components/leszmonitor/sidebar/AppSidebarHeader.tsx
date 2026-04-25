import {
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
} from "@/components/ui/sidebar.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { Link } from "@tanstack/react-router";
import { LeszmonitorLogo } from "@/components/leszmonitor/ui/LeszmonitorLogo.tsx";
import { ProjectSelector } from "@/components/leszmonitor/sidebar/OrgSelector.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import type { Project } from "@/lib/types.ts";
import { WebSocketStatusIndicator } from "@/components/leszmonitor/sidebar/WebSocketStatusIndicator.tsx";

export interface AppSidebarHeaderProps {
  projects: Project[];
}

export function AppSidebarHeader({ projects }: AppSidebarHeaderProps) {
  return (
    <SidebarHeader>
      <Flex direction="row" className="justify-between items-center">
        <div className="p-2">
          <Link to={"/"}>
            <LeszmonitorLogo />
          </Link>
        </div>
        <WebSocketStatusIndicator />
      </Flex>

      <SidebarMenu>
        <SidebarMenuItem>
          {projects ? (
            <ProjectSelector projects={projects} />
          ) : (
            <Skeleton className="h-8 w-full" />
          )}
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>
  );
}
