import { MonitorsApi } from "@/features/monitors/monitors-api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { MonitorForm } from "@/features/monitors/forms/MonitorForm";
import {
  mapMonitorToFormValues,
  type MonitorFormValues,
} from "@/features/monitors/schema";
import { QUERY_KEYS } from "@/lib/consts";
import { useNavigate } from "@tanstack/react-router";

export interface MonitorEditPageProps {
  monitorSlug: string;
}

export function MonitorEditPage({ monitorSlug }: MonitorEditPageProps) {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const { data: monitor } = useQuery({
    queryKey: [QUERY_KEYS.MONITORS, monitorSlug],
    queryFn: () => MonitorsApi.getBySlug(monitorSlug),
  });

  const updateMonitorMutation = useMutation({
    mutationFn: (values: MonitorFormValues) =>
      MonitorsApi.update(monitor?.id ?? monitorSlug, {
        ...values,
        id: monitor?.id ?? monitorSlug,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.MONITORS] });
    },
  });

  const onSubmit = async (value: MonitorFormValues) => {
    await updateMonitorMutation.mutateAsync(value);
    await navigate({
      to: "/monitors/$monitorSlug",
      params: { monitorSlug: value.slug },
    });
  };

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
            onSubmit={onSubmit}
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
