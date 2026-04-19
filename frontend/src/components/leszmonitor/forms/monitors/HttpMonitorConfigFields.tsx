import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { Field, FieldLabel } from "@/components/ui/field.tsx";
import { LMInputField } from "@/components/leszmonitor/forms/inputs/LMInputField.tsx";
import { LMSelect } from "@/components/leszmonitor/forms/inputs/LMSelect.tsx";
import { LMTextareaField } from "@/components/leszmonitor/forms/inputs/LMTextareaField.tsx";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/leszmonitor/forms/inputs/utils.ts";
import type { HttpMonitorFormApi } from "@/components/leszmonitor/forms/NewMonitorForm.tsx";

const httpMethodItems = [
  { value: "GET", label: "GET" },
  { value: "POST", label: "POST" },
  { value: "PUT", label: "PUT" },
  { value: "DELETE", label: "DELETE" },
  { value: "PATCH", label: "PATCH" },
];

export function HttpMonitorConfigFields({
  form,
}: {
  form: HttpMonitorFormApi;
}) {
  return (
    <Flex direction="column" className="gap-4 items-stretch">
      <div className="text-lg font-semibold">HTTP Settings</div>

      <form.Field
        name="config.method"
        children={(field) => (
          <Field>
            <FieldLabel>Method</FieldLabel>
            <LMSelect
              id={field.name}
              name={field.name}
              value={field.state.value}
              onValueChange={(value) =>
                field.handleChange(
                  value as "GET" | "POST" | "PUT" | "DELETE" | "PATCH",
                )
              }
              placeholder="Select HTTP Method"
              items={httpMethodItems}
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      />

      <form.Field
        name="config.url"
        children={(field) => (
          <Field>
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
      />

      <form.Field
        name="config.body"
        children={(field) => (
          <Field>
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
      />

      <form.Field
        name="config.expectedResponseTimeMs"
        children={(field) => (
          <Field>
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
      />

      <form.Field
        name="config.expectedBodyRegex"
        children={(field) => (
          <Field>
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
      />

      {/* TODO: config.headers — needs a key-value pair editor component */}
      {/* TODO: config.expectedHeaders — needs a key-value pair editor component */}
      {/* TODO: config.expectedStatusCodes — needs a multi-number input or tag input */}
      {/* TODO: config.saveResponseBody — needs a checkbox/switch component */}
      {/* TODO: config.saveResponseHeaders — needs a checkbox/switch component */}
    </Flex>
  );
}
