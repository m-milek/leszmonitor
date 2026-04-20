import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { useAtomValue } from "jotai";
import { projectAtom } from "@/lib/atoms.ts";

export const Route = createFileRoute(
  "/_authenticated/projects/$projectId/",
)({
  component: ProjectDashboard,
});

function ProjectDashboard() {
  const { projectId } = Route.useParams();
  const project = useAtomValue(projectAtom);

  return (
    <MainPanelContainer>
      <TypographyH1>Dashboard: {project?.name ?? projectId}</TypographyH1>
    </MainPanelContainer>
  );
}
