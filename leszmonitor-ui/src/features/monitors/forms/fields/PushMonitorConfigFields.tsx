import { TypographyH3 } from "@/components/common/Typography";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { LMInputField } from "@/components/form/LMInputField";
import { getFirstError, isFieldInvalid } from "@/components/form/field-state";
import type { MonitorFormApi } from "@/features/monitors/hooks/useMonitorForm";

export function PushMonitorConfigFields({
  form,
}: Readonly<{ form: MonitorFormApi }>) {
  return (
    <FieldGroup>
      <TypographyH3>Push Monitor Settings</TypographyH3>

      <form.Field name="probeConfig.gracePeriodSeconds">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Grace period (s)</FieldLabel>
            <LMInputField
              name={field.name}
              type="number"
              inputMode="numeric"
              value={field.state.value}
              onChange={(e) => field.handleChange(Number(e.target.value))}
              placeholder="30"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>
    </FieldGroup>
  );
}
