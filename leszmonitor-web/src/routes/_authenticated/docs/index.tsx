import { createFileRoute } from "@tanstack/react-router";
import { PageContainer } from "@/components/leszmonitor/PageContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";

export const Route = createFileRoute("/_authenticated/docs/")({
  component: DocsComponent,
});

function DocsComponent() {
  return (
    <PageContainer>
      <TypographyH1>Documentation</TypographyH1>
    </PageContainer>
  );
}
