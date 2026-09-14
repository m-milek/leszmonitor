import { createFileRoute } from "@tanstack/react-router";
import { PageContainer } from "@/components/leszmonitor/PageContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { Card, CardContent, CardFooter } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { NewMonitorForm } from "@/components/leszmonitor/forms/NewMonitorForm.tsx";

export const Route = createFileRoute("/_authenticated/monitors/new/")({
  component: NewMonitorComponent,
});

function NewMonitorComponent() {
  return (
    <PageContainer>
      <TypographyH1>New Monitor Wizard</TypographyH1>
      <Card>
        <CardContent>
          <NewMonitorForm formId="new-monitor-form" />
        </CardContent>
        <CardFooter>
          <Button type="submit" form="new-monitor-form">
            Create Monitor
          </Button>
        </CardFooter>
      </Card>
    </PageContainer>
  );
}
