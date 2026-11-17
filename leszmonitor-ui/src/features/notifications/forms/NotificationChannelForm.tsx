import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { LMInputField } from "@/components/form/LMInputField";
import { LMSelect } from "@/components/form/LMSelect";
import { LMSwitch } from "@/components/form/LMSwitch";
import { getFirstError, isFieldInvalid } from "@/components/form/field-state";
import {
  defaultChannelConfigs,
  newNotificationChannelDefaultValues,
  type NotificationChannelFormValues,
} from "@/features/notifications/schema";
import {
  notificationChannelKindLabels,
  notificationChannelKinds,
  type NotificationChannelKind,
} from "@/features/notifications/types";
import { useNotificationChannelForm } from "@/features/notifications/hooks/useNotificationChannelForm";
import { WebhookConfigFields } from "@/features/notifications/forms/fields/WebhookConfigFields";
import { NtfyConfigFields } from "@/features/notifications/forms/fields/NtfyConfigFields";
import { GotifyConfigFields } from "@/features/notifications/forms/fields/GotifyConfigFields";

const kindItems = notificationChannelKinds.map((kind) => ({
  value: kind,
  label: notificationChannelKindLabels[kind],
}));

const isChannelKind = (value: string): value is NotificationChannelKind =>
  (notificationChannelKinds as readonly string[]).includes(value);

export interface NotificationChannelFormProps {
  id?: string;
  defaultValues?: NotificationChannelFormValues;
  isEdit?: boolean;
  onSubmit: (values: NotificationChannelFormValues) => Promise<void>;
}

export function NotificationChannelForm({
  id = "notification-channel-form",
  defaultValues = newNotificationChannelDefaultValues,
  isEdit = false,
  onSubmit,
}: Readonly<NotificationChannelFormProps>) {
  const form = useNotificationChannelForm({
    defaultValues,
    onSubmit,
    isEdit,
  });

  return (
    <form
      id={id}
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}
    >
      <FieldGroup>
        <form.Field
          name="kind"
          listeners={{
            onChange: ({ value }) => {
              form.setFieldValue("config", defaultChannelConfigs[value]);
            },
          }}
        >
          {(field) => (
            <Field id={field.name}>
              <FieldLabel>Type</FieldLabel>
              <LMSelect
                id={field.name}
                name={field.name}
                value={field.state.value}
                onValueChange={(value) => {
                  if (isChannelKind(value)) {
                    field.handleChange(value);
                  }
                }}
                items={kindItems}
                disabled={isEdit}
                isInvalid={isFieldInvalid(field)}
                errorMessage={getFirstError(field)}
              />
            </Field>
          )}
        </form.Field>

        <form.Field name="name">
          {(field) => (
            <Field id={field.name}>
              <FieldLabel>Channel Name</FieldLabel>
              <LMInputField
                name={field.name}
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                isInvalid={isFieldInvalid(field)}
                errorMessage={getFirstError(field)}
              />
            </Field>
          )}
        </form.Field>

        <form.Field name="enabled">
          {(field) => (
            <Field id={field.name} orientation="horizontal">
              <LMSwitch
                name={field.name}
                checked={field.state.value}
                onCheckedChange={field.handleChange}
              />
              <FieldLabel>Enabled</FieldLabel>
            </Field>
          )}
        </form.Field>

        <form.Subscribe selector={(state) => state.values.kind}>
          {(kind) => {
            switch (kind) {
              case "webhook":
                return <WebhookConfigFields form={form} isEdit={isEdit} />;
              case "ntfy":
                return <NtfyConfigFields form={form} isEdit={isEdit} />;
              case "gotify":
                return <GotifyConfigFields form={form} isEdit={isEdit} />;
              default:
                return null;
            }
          }}
        </form.Subscribe>
      </FieldGroup>
    </form>
  );
}
