import { createFileRoute } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { Card, CardContent, CardFooter } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { getMonitorBySlug, updateMonitor } from "@/lib/data/monitorData.ts";
import { MonitorForm } from "@/components/leszmonitor/forms/NewMonitorForm.tsx";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { mapMonitorToFormValues, type MonitorFormValues } from "@/lib/types.ts";

export const Route = createFileRoute(
  "/_authenticated/projects/$projectId/monitors/$monitorSlug/edit",
)({
  component: MonitorEditRoute,
});

function MonitorEditRoute() {
  const { projectId, monitorSlug } = Route.useParams();
  const queryClient = useQueryClient();

  const { data: monitor } = useQuery({
    queryKey: [QUERY_KEYS.MONITORS, monitorSlug, projectId],
    queryFn: () => getMonitorBySlug(projectId, monitorSlug),
  });

  const updateMonitorMutation = useMutation({
    mutationFn: (values: MonitorFormValues) =>
      updateMonitor(monitor?.id ?? monitorSlug, {
        ...values,
        id: monitor?.id ?? monitorSlug,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.MONITORS] });
    },
  });

  if (!monitor) {
    return null;
  }

  return (
    <MainPanelContainer>
      <TypographyH1>Edit Monitor</TypographyH1>
      <Card>
        <CardContent>
          <MonitorForm
            formId="edit-monitor-form"
            projectSlug={projectId}
            defaultValues={mapMonitorToFormValues(monitor)}
            onSubmit={(value) => updateMonitorMutation.mutateAsync(value)}
          />
        </CardContent>
        <CardFooter>
          <Button type="submit" form="edit-monitor-form">
            Save Changes
          </Button>
        </CardFooter>
      </Card>
    </MainPanelContainer>
  );
}
