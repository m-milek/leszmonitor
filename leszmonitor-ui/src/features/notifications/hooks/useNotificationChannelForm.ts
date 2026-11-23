import { useForm } from "@tanstack/react-form";
import {
  type NotificationChannelFormValues,
  buildNotificationChannelSchema,
} from "@/features/notifications/schema";

interface UseNotificationChannelFormOptions {
  defaultValues: NotificationChannelFormValues;
  onSubmit: (value: NotificationChannelFormValues) => Promise<void>;
  isEdit: boolean;
}

export function useNotificationChannelForm({
  defaultValues,
  onSubmit,
  isEdit,
}: UseNotificationChannelFormOptions) {
  const schema = buildNotificationChannelSchema(isEdit);
  return useForm({
    defaultValues,
    validators: {
      onSubmit: schema,
    },
    onSubmit: async ({ value }) => {
      await onSubmit(schema.parse(value));
    },
  });
}

export type NotificationChannelFormApi = ReturnType<
  typeof useNotificationChannelForm
>;
