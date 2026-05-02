
import type { MonitorFormValues } from "./types.ts";

import { useForm, type FormValidateOrFn } from "@tanstack/react-form";
import { newMonitorSchema } from "./types.ts";

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function __monitorFormHelper() {
  return useForm({
    defaultValues: {} as MonitorFormValues,
    validators: {
      onSubmit:
        newMonitorSchema as unknown as FormValidateOrFn<MonitorFormValues>,
    },
  });
}
export type MonitorFormApi = ReturnType<typeof __monitorFormHelper>;



