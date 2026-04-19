import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { Card, CardContent, CardFooter } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useQuery } from "@tanstack/react-query";
import { getProjects } from "@/lib/data/projectData.ts";
import { NewMonitorForm } from "@/components/leszmonitor/forms/NewMonitorForm.tsx";

export const Route = createFileRoute(
  "/_authenticated/org/$orgId/monitors/new/",
)({
  component: NewMonitorComponent,
});

function NewMonitorComponent() {
  const { orgId } = Route.useParams();

  const { data: projects } = useQuery({
    queryKey: ["projects", orgId],
    queryFn: () => getProjects(orgId),
  });

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
            projects={projects}
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
