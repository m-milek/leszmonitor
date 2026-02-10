import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/sidebar/Typography.tsx";

export const Route = createFileRoute("/_authenticated/$teamId/dashboard/")({
  component: DashboardComponent,
});

function DashboardComponent() {
  return (
    <MainPanelContainer>
      <TypographyH1>Dashboard</TypographyH1>
    </MainPanelContainer>
  );
}
