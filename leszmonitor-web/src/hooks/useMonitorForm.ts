import { useForm } from "@tanstack/react-form";
import { type MonitorFormValues, newMonitorSchema } from "@/lib/types.ts";

interface UseMonitorFormOptions {
  defaultValues: MonitorFormValues;
  onSubmit: (value: MonitorFormValues) => Promise<void>;
  onReset?: () => void;
}

export function useMonitorForm({
  defaultValues,
  onSubmit,
  onReset,
}: UseMonitorFormOptions) {
  return useForm({
    defaultValues,
    validators: {
      onSubmit: ({ value }) => {
        const result = newMonitorSchema.safeParse(value);
        if (!result.success) {
          return {
            formErrors: result.error.issues.map((i) => i.message),
          };
        }
        return undefined;
      },
    },
    onSubmit: async ({ value }) => {
      await onSubmit(value);
      onReset?.();
    },
    onSubmitInvalid: () => {
      console.log("Invalid form submission");
    },
  });
}

export type MonitorFormApi = ReturnType<typeof useMonitorForm>;
