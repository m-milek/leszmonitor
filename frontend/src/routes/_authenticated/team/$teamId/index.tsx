import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/Typography.tsx";

export const Route = createFileRoute("/_authenticated/team/$teamId/")({
  component: TeamDashboard,
});

function TeamDashboard() {
  const { teamId } = Route.useParams();

  return (
    <MainPanelContainer>
      <TypographyH1>Team Dashboard for Team ID: {teamId}</TypographyH1>
    </MainPanelContainer>
  );
}
