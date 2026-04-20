import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { Card, CardContent, CardFooter } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { NewMonitorForm } from "@/components/leszmonitor/forms/NewMonitorForm.tsx";

export const Route = createFileRoute(
  "/_authenticated/projects/$projectId/monitors/new/",
)({
  component: NewMonitorComponent,
});

function NewMonitorComponent() {
  const { projectId } = Route.useParams();

  const onSubmit = async (value: unknown) => {
    console.log(value);
  };

  return (
    <MainPanelContainer>
      <TypographyH1>New Monitor Wizard</TypographyH1>
      <Card>
        <CardContent>
          <NewMonitorForm
            formId="new-monitor-form"
            onSubmitMonitor={onSubmit}
            projectId={projectId}
          />
        </CardContent>
        <CardFooter>
          <Button type="submit" form="new-monitor-form">
            Create Monitor
          </Button>
        </CardFooter>
      </Card>
    </MainPanelContainer>
  );
}
