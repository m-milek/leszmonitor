import type { Project } from "@/lib/types.ts";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { colorFromString } from "@/lib/colorFromString.ts";
import { Card } from "@/components/ui/card.tsx";
import { LucideHeartPulse, LucideUser } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { getMonitorsByProjectSlug } from "@/lib/data/monitorData.ts";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import { useNavigate } from "@tanstack/react-router";
import { StatusDot } from "@/components/leszmonitor/ui/StatusDot.tsx";

export interface ProjectListItemProps {
  project: Project;
}

export const ProjectListItem = ({ project }: ProjectListItemProps) => {
  const backgroundColor = colorFromString(project.id);

  const { data: monitors } = useQuery({
    queryKey: ["monitors", project.slug],
    queryFn: async () => getMonitorsByProjectSlug(project.slug),
  });

  const navigate = useNavigate();

  const onClick = () => {
    navigate({
      to: "/projects/$projectId",
      params: { projectId: project.slug },
    });
  };

  const ProjectIcon = () => (
    <Flex
      className="items-center justify-center min-w-12 max-w-12 min-h-12 max-h-12 rounded-full border border-border text-slate-800 shrink-0"
      style={{ backgroundColor }}
    >
      <span className="text-2xl">
        {project.name.substring(0, 1).toUpperCase()}
      </span>
    </Flex>
  );

  const MemberCount = () => (
    <Flex className="items-center gap-1 text-muted-foreground text-sm">
      <LucideUser /> {project.members.length} members
    </Flex>
  );

  const MonitorCount = () => (
    <Flex className="items-center gap-1 text-muted-foreground text-sm">
      {monitors ? (
        <>
          <LucideHeartPulse /> {monitors ? monitors.length : 0} monitors
        </>
      ) : (
        <Skeleton className="h-4 w-8" />
      )}
    </Flex>
  );

  return (
    <Card
      className="min-w-0 cursor-pointer hover:bg-secondary"
      onClick={onClick}
    >
      <Flex className="items-center justify-between p-4 border-border gap-4 w-full min-w-0">
        <Flex className="items-center flex-1 min-w-0 gap-4">
          <ProjectIcon />

          <Flex direction="column" className="flex-1 min-w-0 w-0">
            <Flex className="text-lg gap-2 font-medium truncate items-center">
              <span>{project.name}</span>
              <span className="text-sm text-muted-foreground">
                {project.slug}
              </span>
            </Flex>

            <div className="text-sm text-muted-foreground truncate">
              {project.description}
            </div>
          </Flex>
        </Flex>
        <MonitorCount />
        <MemberCount />
        <StatusDot status="pending" />
      </Flex>
    </Card>
  );
};
