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
import { LMSwitch } from "@/components/leszmonitor/forms/inputs/LMSwitch.tsx";
import { LMKeyValueInput } from "@/components/leszmonitor/forms/inputs/LMKeyValue.tsx";
import { LMMultiSelect } from "@/components/leszmonitor/forms/inputs/LMMultiSelect.tsx";

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
}: {
  form: HttpMonitorFormApi;
}) {
  return (
    <Flex direction="column" className="gap-4 items-stretch">
      <div className="text-lg font-semibold">HTTP Settings</div>

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
        name="config.headers"
        children={(field) => (
          <Field>
            <FieldLabel>Request Headers</FieldLabel>
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
      />

      <div className="text-lg font-semibold">Expected Response</div>

      <form.Field
        name="config.expectedStatusCodes"
        children={(field) => (
          <Field>
            <FieldLabel>Expected Status Codes</FieldLabel>
            <LMMultiSelect
              name={field.name}
              options={statusCodes}
              value={Array.isArray(field.state.value) ? field.state.value.map(String) : []}
              onChange={(values) => field.handleChange(values.map(Number))}
              placeholder="Add status code"
              emptyMessage="No status codes found."
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

      <form.Field
        name="config.expectedHeaders"
        children={(field) => (
          <Field>
            <FieldLabel>Expected Response Headers</FieldLabel>
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
      />

      <div className="text-lg font-semibold">Capture Response</div>

      <form.Field
        name="config.saveResponseBody"
        children={(field) => (
          <Field>
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
      />
      <form.Field
        name="config.saveResponseHeaders"
        children={(field) => (
          <Field>
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
      />
    </Flex>
  );
}
