import type { Project } from "@/lib/types.ts";
import { SidebarMenu } from "@/components/ui/sidebar.tsx";
import {
  LayoutDashboardIcon,
  LucideActivity,
  LucideLogs,
  LucideNotebookText,
  LucideUsers,
} from "lucide-react";
import { SidebarButton } from "@/components/leszmonitor/sidebar/AppSidebar.tsx";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { useNavigate, useRouterState } from "@tanstack/react-router";
import { useAppStore } from "@/lib/store.ts";

export interface ProjectMenuProps {
  projects: Project[];
}

export const ProjectMenu = ({ projects }: ProjectMenuProps) => {
  const navigate = useNavigate();
  const routerState = useRouterState();
  const { project, setProject } = useAppStore();

  const onProjectSelect = async (projectSlug: string) => {
    if (!projectSlug) {
      setProject(null);
      await navigate({ to: "/projects" });
      return;
    }

    const isInProjectRoute = projects.some((p) =>
      routerState.location.pathname.startsWith(`/projects/${p.slug}`),
    );
    if (isInProjectRoute) {
      // navigate to the same location but in a different project
      const newPath = routerState.location.pathname.replace(
        /\/projects\/[^/]+/,
        `/projects/${projectSlug}`,
      );
      await navigate({
        to: newPath,
      });
      return;
    }

    await navigate({
      to: "/projects/$projectId",
      params: { projectId: projectSlug },
    });
  };

  return (
    <Flex
      className="border border-sidebar-border rounded-sm p-2 gap-2 bg-secondary"
      direction="column"
    >
      <Select value={project?.slug ?? ""} onValueChange={onProjectSelect}>
        <SelectTrigger className="w-full" withClearButton>
          <SelectValue placeholder="Select project" />
        </SelectTrigger>
        <SelectContent className="w-full" position="popper">
          {projects.map((project) => (
            <SelectItem value={project.slug} key={project.slug}>
              {project.name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      {project && (
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
            icon={<LucideNotebookText />}
            href={`/projects/${project.slug}/events`}
            label="Events"
          />
          <SidebarButton
            icon={<LucideUsers />}
            href={`/projects/${project.slug}/members`}
            label="Access"
          />
          <SidebarButton
            icon={<LucideLogs />}
            href={`/projects/${project.slug}/audit-log`}
            label="Audit Log"
          />
        </SidebarMenu>
      )}
    </Flex>
  );
};
