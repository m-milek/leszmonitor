import { useForm } from "@tanstack/react-form";
import {
  type NotificationChannelFormValues,
  notificationChannelSchema,
} from "@/features/notifications/schema";

interface UseNotificationChannelFormOptions {
  defaultValues: NotificationChannelFormValues;
  onSubmit: (value: NotificationChannelFormValues) => Promise<void>;
}

export function useNotificationChannelForm({
  defaultValues,
  onSubmit,
}: UseNotificationChannelFormOptions) {
  return useForm({
    defaultValues,
    validators: {
      onSubmit: notificationChannelSchema,
    },
    onSubmit: async ({ value }) => {
      await onSubmit(notificationChannelSchema.parse(value));
    },
  });
}

export type NotificationChannelFormApi = ReturnType<
  typeof useNotificationChannelForm
>;
