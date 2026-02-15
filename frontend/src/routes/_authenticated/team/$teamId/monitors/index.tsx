import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/Typography.tsx";

export const Route = createFileRoute("/_authenticated/team/$teamId/monitors/")({
  component: MonitorsComponent,
});

function MonitorsComponent() {
  return (
    <MainPanelContainer>
      <TypographyH1>Monitors</TypographyH1>
    </MainPanelContainer>
  );
}
