import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";

export const Route = createFileRoute("/_authenticated/org/$orgId/")({
  component: OrgDashboard,
});

function OrgDashboard() {
  const { orgId } = Route.useParams();

  return (
    <MainPanelContainer>
      <TypographyH1>Org Dashboard for Org ID: {orgId}</TypographyH1>
    </MainPanelContainer>
  );
}
