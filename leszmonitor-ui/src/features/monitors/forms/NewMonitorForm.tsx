import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { createMonitor } from "@/features/monitors/monitors-api";
import { QUERY_KEYS } from "@/lib/consts";
import type { MonitorFormValues } from "@/features/monitors/schema";
import { MonitorForm } from "@/features/monitors/forms/MonitorForm";

export interface NewMonitorFormProps {
  formId?: string;
}

export function NewMonitorForm({
  formId = "new-monitor-form",
}: Readonly<NewMonitorFormProps>) {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const createMonitorMutation = useMutation({
    mutationFn: (monitor: MonitorFormValues) => createMonitor(monitor),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.MONITORS] });
    },
  });

  const onSubmit = async (value: MonitorFormValues) => {
    await createMonitorMutation.mutateAsync(value);
    await navigate({
      to: "/monitors/$monitorSlug",
      params: { monitorSlug: value.slug },
    });
  };

  return <MonitorForm formId={formId} onSubmit={onSubmit} resetOnSuccess />;
}
