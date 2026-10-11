import { TypographyH3 } from "@/components/common/Typography";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { LMInputField } from "@/components/form/LMInputField";
import { getFirstError, isFieldInvalid } from "@/components/form/field-state";
import type { MonitorFormApi } from "@/features/monitors/hooks/useMonitorForm";

export function PingMonitorConfigFields({
  form,
}: Readonly<{ form: MonitorFormApi }>) {
  return (
    <FieldGroup>
      <TypographyH3>Ping Monitor Settings</TypographyH3>

      <form.Field name="probeConfig.host">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Host</FieldLabel>
            <LMInputField
              name={field.name}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              placeholder="192.168.1.1"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="probeConfig.count">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Packets</FieldLabel>
            <LMInputField
              name={field.name}
              type="number"
              inputMode="numeric"
              value={field.state.value}
              onChange={(e) => field.handleChange(Number(e.target.value))}
              placeholder="3"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="probeConfig.timeoutMs">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Timeout (ms)</FieldLabel>
            <LMInputField
              name={field.name}
              type="number"
              inputMode="numeric"
              value={field.state.value}
              onChange={(e) => field.handleChange(Number(e.target.value))}
              placeholder="5000"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>
    </FieldGroup>
  );
}
