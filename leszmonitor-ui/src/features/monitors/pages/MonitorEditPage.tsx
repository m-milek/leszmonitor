import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { PageContainer } from "@/components/common/PageContainer.tsx";
import { TypographyH1 } from "@/components/common/Typography.tsx";
import { Card, CardContent, CardFooter } from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
  getMonitorBySlug,
  updateMonitor,
} from "@/features/monitors/api/monitors.ts";
import { MonitorForm } from "@/features/monitors/forms/MonitorForm.tsx";
import {
  mapMonitorToFormValues,
  type MonitorFormValues,
} from "@/features/monitors/model/schema.ts";
import { QUERY_KEYS } from "@/lib/consts.ts";

export interface MonitorEditPageProps {
  monitorSlug: string;
}

export function MonitorEditPage({ monitorSlug }: MonitorEditPageProps) {
  const queryClient = useQueryClient();

  const { data: monitor } = useQuery({
    queryKey: [QUERY_KEYS.MONITORS, monitorSlug],
    queryFn: () => getMonitorBySlug(monitorSlug),
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
    <PageContainer>
      <TypographyH1>Edit Monitor</TypographyH1>
      <Card>
        <CardContent>
          <MonitorForm
            formId="edit-monitor-form"
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
    </PageContainer>
  );
}
