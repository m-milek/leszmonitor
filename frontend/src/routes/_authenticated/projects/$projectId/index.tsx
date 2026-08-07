import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import {
  TypographyH1,
  TypographyH2,
} from "@/components/leszmonitor/ui/Typography.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex";
import { Initial } from "@/components/leszmonitor/Initial.tsx";
import { getProjectBySlug } from "@/lib/data/projectData.ts";
import { useQuery } from "@tanstack/react-query";
import { getMonitorsByProjectSlug } from "@/lib/data/monitorData.ts";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { MonitorListItem } from "@/components/leszmonitor/MonitorListItem.tsx";

export const Route = createFileRoute("/_authenticated/projects/$projectId/")({
  component: ProjectDashboard,
  loader: async ({ params }) => getProjectBySlug(params.projectId),
});

function ProjectDashboard() {
  const project = Route.useLoaderData();

  const { data: monitors } = useQuery({
    queryFn: async () => getMonitorsByProjectSlug(project.slug),
    queryKey: [QUERY_KEYS.MONITORS, project.slug],
  });

  return (
    <MainPanelContainer>
      <Flex direction="column" className="gap-4">
        <Flex className="gap-4">
          <Initial
            text={project.name}
            textForColorCalculation={project.id}
            size="xl"
          />
          <div className="flex flex-col justify-center">
            <TypographyH1>{project.name}</TypographyH1>
            <span className="text-muted-foreground">{project.id}</span>
          </div>
        </Flex>
        <Flex className="gap-4">
          <Card className="flex-1">
            <CardHeader>
              <TypographyH2>Monitors</TypographyH2>
            </CardHeader>
            <CardContent>
              <Flex direction="column" className="gap-4">
                {monitors?.map((monitor) => (
                  <MonitorListItem
                    key={monitor.id}
                    monitor={monitor}
                    projectSlug={project.slug}
                  />
                ))}
              </Flex>
            </CardContent>
          </Card>
          <Card className="flex-1">
            <CardHeader>
              <TypographyH2>Members</TypographyH2>
            </CardHeader>
            <CardContent>
              <Flex direction="column" className="gap-4">
                {project.members?.map((monitor) => (
                  <div key={monitor.id} className="flex items-center gap-2">
                    <Initial
                      text={monitor.username}
                      textForColorCalculation={monitor.id}
                      size="md"
                    />
                    <span>{monitor.username}</span>
                  </div>
                ))}
              </Flex>
            </CardContent>
          </Card>
        </Flex>
      </Flex>
    </MainPanelContainer>
  );
}
