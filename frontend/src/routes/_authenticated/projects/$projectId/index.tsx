import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { useAppStore } from "@/lib/store.ts";

export const Route = createFileRoute(
  "/_authenticated/projects/$projectId/",
)({
  component: ProjectDashboard,
});

function ProjectDashboard() {
  const { projectId } = Route.useParams();
  const { project } = useAppStore();

  return (
    <MainPanelContainer>
      <TypographyH1>Dashboard: {project?.name ?? projectId}</TypographyH1>
    </MainPanelContainer>
  );
}
