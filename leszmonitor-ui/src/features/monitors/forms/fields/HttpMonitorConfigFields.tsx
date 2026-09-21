import { TypographyH3 } from "@/components/common/Typography";
import { Flex } from "@/components/common/Flex";
import {
  Field,
  FieldGroup,
  FieldLabel,
  FieldTitle,
} from "@/components/ui/field";
import { LMInputField } from "@/components/form/LMInputField";
import { LMSelect } from "@/components/form/LMSelect";
import { LMTextareaField } from "@/components/form/LMTextareaField";
import { getFirstError, isFieldInvalid } from "@/components/form/field-state";
import { LMSwitch } from "@/components/form/LMSwitch";
import { LMKeyValueInput } from "@/components/form/LMKeyValue";
import { LMMultiSelect } from "@/components/form/LMMultiSelect";
import type { HttpMethod } from "@/features/monitors/types";
import type { MonitorFormApi } from "@/features/monitors/hooks/useMonitorForm";

const httpMethodItems = [
  { value: "GET", label: "GET" },
  { value: "POST", label: "POST" },
  { value: "PUT", label: "PUT" },
  { value: "DELETE", label: "DELETE" },
  { value: "PATCH", label: "PATCH" },
];

const statusCodes = Array.from({ length: 600 }, (_, i) => i)
  .filter((code) => code >= 100 && code < 600)
  .map((code) => String(code));

export function HttpMonitorConfigFields({
  form,
}: Readonly<{ form: MonitorFormApi }>) {
  return (
    <FieldGroup>
      <TypographyH3>HTTP Settings</TypographyH3>

      <form.Field name="probeConfig.url">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>URL</FieldLabel>
            <LMInputField
              name={field.name}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              placeholder="https://example.com/health"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="probeConfig.method">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Method</FieldLabel>
            <LMSelect
              id={field.name}
              name={field.name}
              value={field.state.value}
              onValueChange={(value) => field.handleChange(value as HttpMethod)}
              placeholder="Select HTTP Method"
              items={httpMethodItems}
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="probeConfig.body">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Request Body</FieldLabel>
            <LMTextareaField
              name={field.name}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              placeholder='{"key": "value"}'
              rows={4}
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="probeConfig.headers">
        {(field) => (
          <Field id={field.name}>
            <FieldTitle>Request Headers</FieldTitle>
            <LMKeyValueInput
              name={field.name}
              value={field.state.value}
              onChange={(value) => field.handleChange(value)}
              keyPlaceholder="Header"
              valuePlaceholder="Header Value"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <TypographyH3>Expected Response</TypographyH3>

      <form.Field name="probeConfig.expectedStatusCodes">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Expected Status Codes</FieldLabel>
            <LMMultiSelect
              id={field.name}
              name={field.name}
              options={statusCodes}
              value={
                Array.isArray(field.state.value)
                  ? field.state.value.map(String)
                  : []
              }
              onChange={(values) => field.handleChange(values.map(Number))}
              placeholder="Add status code"
              emptyMessage="No status codes found."
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="probeConfig.expectedResponseTimeMs">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Expected Response Time (ms)</FieldLabel>
            <LMInputField
              name={field.name}
              type="number"
              inputMode="numeric"
              value={field.state.value ?? ""}
              onChange={(e) =>
                field.handleChange(
                  e.target.value ? Number(e.target.value) : undefined,
                )
              }
              placeholder="1000"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="probeConfig.expectedBodyRegex">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Expected Body Pattern (RegExp)</FieldLabel>
            <LMInputField
              name={field.name}
              value={field.state.value ?? ""}
              onChange={(e) => field.handleChange(e.target.value || undefined)}
              placeholder="OK|healthy"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="probeConfig.expectedHeaders">
        {(field) => (
          <Field id={field.name}>
            <FieldTitle>Expected Response Headers</FieldTitle>
            <LMKeyValueInput
              name={field.name}
              value={field.state.value}
              onChange={(value) => field.handleChange(value)}
              keyPlaceholder="Header"
              valuePlaceholder="Header Value"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <TypographyH3>Capture Response</TypographyH3>

      <form.Field name="probeConfig.saveResponseBody">
        {(field) => (
          <Field id={field.name}>
            <Flex direction="row" className="justify-between">
              <FieldLabel>Save Response Body</FieldLabel>
              <LMSwitch
                name={field.name}
                checked={!!field.state.value}
                onCheckedChange={(checked) =>
                  field.handleChange(checked ? true : undefined)
                }
              />
            </Flex>
          </Field>
        )}
      </form.Field>
      <form.Field name="probeConfig.saveResponseHeaders">
        {(field) => (
          <Field id={field.name}>
            <Flex direction="row" className="justify-between">
              <FieldLabel>Save Response Headers</FieldLabel>
              <LMSwitch
                name={field.name}
                checked={!!field.state.value}
                onCheckedChange={(checked) =>
                  field.handleChange(checked ? true : undefined)
                }
              />
            </Flex>
          </Field>
        )}
      </form.Field>
    </FieldGroup>
  );
}
