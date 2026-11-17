import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { LMInputField } from "@/components/form/LMInputField";
import { getFirstError, isFieldInvalid } from "@/components/form/field-state";
import type { NotificationChannelFormApi } from "@/features/notifications/hooks/useNotificationChannelForm";

export interface GotifyConfigFieldsProps {
  form: NotificationChannelFormApi;
  isEdit: boolean;
}

export function GotifyConfigFields({
  form,
  isEdit,
}: Readonly<GotifyConfigFieldsProps>) {
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
              placeholder="https://gotify.example.com"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="config.appToken">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>App token</FieldLabel>
            <LMInputField
              name={field.name}
              type="password"
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              placeholder={isEdit ? "Unchanged" : undefined}
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="config.priority">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Priority (0-10)</FieldLabel>
            <LMInputField
              name={field.name}
              type="number"
              inputMode="numeric"
              value={field.state.value}
              onChange={(e) => field.handleChange(Number(e.target.value))}
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>
    </FieldGroup>
  );
}
