import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { LMInputField } from "@/components/form/LMInputField";
import { LMSelect } from "@/components/form/LMSelect";
import { getFirstError, isFieldInvalid } from "@/components/form/field-state";
import type { NotificationChannelFormApi } from "@/features/notifications/hooks/useNotificationChannelForm";

const priorityItems = [
  { value: "0", label: "Default" },
  { value: "1", label: "1 (min)" },
  { value: "2", label: "2 (low)" },
  { value: "3", label: "3 (default)" },
  { value: "4", label: "4 (high)" },
  { value: "5", label: "5 (max)" },
];

export interface NtfyConfigFieldsProps {
  form: NotificationChannelFormApi;
}

export function NtfyConfigFields({ form }: Readonly<NtfyConfigFieldsProps>) {
  return (
    <FieldGroup>
      <form.Field name="config.serverUrl">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Server URL</FieldLabel>
            <LMInputField
              name={field.name}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              placeholder="https://ntfy.sh"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="config.topic">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Topic</FieldLabel>
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

      <form.Field name="config.token">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Access token</FieldLabel>
            <LMInputField
              name={field.name}
              type="password"
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              placeholder="Optional"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="config.username">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Username</FieldLabel>
            <LMInputField
              name={field.name}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              placeholder="Optional"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="config.password">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Password</FieldLabel>
            <LMInputField
              name={field.name}
              type="password"
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              placeholder="Optional"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="config.priority">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Priority</FieldLabel>
            <LMSelect
              id={field.name}
              name={field.name}
              value={String(field.state.value)}
              onValueChange={(value) => field.handleChange(Number(value))}
              items={priorityItems}
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>
    </FieldGroup>
  );
}
