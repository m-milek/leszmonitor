import {
  type MonitorFormValues,
  defaultConfigs,
  newMonitorSchemaDefaultValues,
} from "@/features/monitors/schema";

export const buildMonitorDefaults = (
  defaultValues?: Partial<MonitorFormValues>,
): MonitorFormValues => {
  const baseValues = {
    ...newMonitorSchemaDefaultValues,
  };
  const type = defaultValues?.type;

  if (!type) {
    return {
      ...baseValues,
      ...defaultValues,
    } as MonitorFormValues;
  }

  return {
    ...baseValues,
    ...defaultValues,
    type,
    probeConfig: {
      ...defaultConfigs[type],
      ...(defaultValues?.probeConfig ?? {}),
    },
  } as MonitorFormValues;
};
