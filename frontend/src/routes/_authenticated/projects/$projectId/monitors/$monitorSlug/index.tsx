import { createFileRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { getMonitorBySlug } from "@/lib/data/monitorData.ts";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";

export const Route = createFileRoute(
  "/_authenticated/projects/$projectId/monitors/$monitorSlug/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const { projectId, monitorSlug } = Route.useParams();

  const { data: monitor } = useQuery({
    queryKey: ["monitors", monitorSlug, projectId],
    queryFn: () => getMonitorBySlug(projectId, monitorSlug),
  });

  if (!monitor) {
    return null;
  }

  return (
    <MainPanelContainer>
      <TypographyH1>{monitor.name}</TypographyH1>
      <span>{monitor.id}</span>
      <span>{monitor.type}</span>
      <span>{monitor.description}</span>
    </MainPanelContainer>
  );
}
