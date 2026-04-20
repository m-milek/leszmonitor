import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";

export const Route = createFileRoute(
  "/_authenticated/projects/$projectId/dashboard/",
)({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <MainPanelContainer>
      <TypographyH1>Project Dashboard</TypographyH1>
    </MainPanelContainer>
  );
}
