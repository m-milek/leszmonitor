import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { LMInputField } from "@/components/form/LMInputField";
import { LMKeyValueInput } from "@/components/form/LMKeyValue";
import { getFirstError, isFieldInvalid } from "@/components/form/field-state";
import type { NotificationChannelFormApi } from "@/features/notifications/hooks/useNotificationChannelForm";

export interface WebhookConfigFieldsProps {
  form: NotificationChannelFormApi;
}

export function WebhookConfigFields({
  form,
}: Readonly<WebhookConfigFieldsProps>) {
  return (
    <FieldGroup>
      <form.Field name="config.url">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>URL</FieldLabel>
            <LMInputField
              name={field.name}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              placeholder="https://example.com/hooks/leszmonitor"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="config.headers">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Headers</FieldLabel>
            <LMKeyValueInput
              name={field.name}
              value={field.state.value}
              onChange={(value) => field.handleChange(value)}
              keyPlaceholder="Authorization"
              valuePlaceholder="Bearer ..."
              addButtonText="Add header"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>
    </FieldGroup>
  );
}
