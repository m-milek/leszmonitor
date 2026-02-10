import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/sidebar/Typography.tsx";

export const Route = createFileRoute("/_authenticated/docs/")({
  component: DocsComponent,
});

function DocsComponent() {
  return (
    <MainPanelContainer>
      <TypographyH1>Documentation</TypographyH1>
    </MainPanelContainer>
  );
}
