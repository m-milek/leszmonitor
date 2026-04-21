import { useForm, type FormValidateOrFn } from "@tanstack/react-form";
import { useState } from "react";
import { slugFromString } from "@/lib/slugFromString.ts";
import {
  defaultConfigs,
  isValidMonitorType,
  type MonitorType,
  newMonitorSchema,
  newMonitorSchemaDefaultValues,
} from "@/lib/types.ts";
import { Field, FieldLabel } from "@/components/ui/field.tsx";
import { LMInputField } from "@/components/leszmonitor/forms/inputs/LMInputField.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { LMSelect } from "@/components/leszmonitor/forms/inputs/LMSelect.tsx";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/leszmonitor/forms/inputs/utils.ts";
import { LMTextareaField } from "@/components/leszmonitor/forms/inputs/LMTextareaField.tsx";
import { Divider } from "@/components/leszmonitor/ui/Divider.tsx";
import { MonitorConfigFields } from "@/components/leszmonitor/forms/monitors/MonitorConfigFields.tsx";
import type {
  HttpMonitorFormValues,
  PingMonitorFormValues,
  MonitorFormValues,
} from "@/lib/types.ts";
import { Switch } from "@/components/ui/switch.tsx";

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function __monitorFormHelper() {
  return useForm({
    defaultValues: {} as MonitorFormValues,
    validators: {
      onSubmit:
        newMonitorSchema as unknown as FormValidateOrFn<MonitorFormValues>,
    },
  });
}
export type MonitorFormApi = ReturnType<typeof __monitorFormHelper>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function __httpMonitorFormHelper() {
  return useForm({
    defaultValues: {} as HttpMonitorFormValues,
  });
}
export type HttpMonitorFormApi = ReturnType<typeof __httpMonitorFormHelper>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function __pingMonitorFormHelper() {
  return useForm({
    defaultValues: {} as PingMonitorFormValues,
  });
}
export type PingMonitorFormApi = ReturnType<typeof __pingMonitorFormHelper>;

export interface NewMonitorFormProps {
  onSubmitMonitor: (value: MonitorFormValues) => Promise<unknown>;
  projectId: string;
  formId?: string;
}

export function NewMonitorForm({
  onSubmitMonitor,
  projectId,
  formId = "new-monitor-form",
}: NewMonitorFormProps) {
  const [isSlugModified, setIsSlugModified] = useState(false);
  const form = useForm({
    defaultValues: {
      ...newMonitorSchemaDefaultValues,
      projectId,
    } as MonitorFormValues,
    validators: {
      onSubmit:
        newMonitorSchema as unknown as FormValidateOrFn<MonitorFormValues>,
    },
    onSubmit: async ({ value }) => {
      await onSubmitMonitor(value as MonitorFormValues);
    },
    onSubmitInvalid: ({ value }) => {
      console.log("Invalid form submission");
      console.log("Values:", value);
    },
  });

  const [useCustomSlug, setUseCustomSlug] = useState(false);

  const onUseCustomSlugChanged = (checked: boolean) => {
    setUseCustomSlug(checked);
    if (!checked) {
      const name = form.state.values.name;
      form.setFieldValue("slug", slugFromString(name));
    }
  };

  const monitorTypeSelectItems: { value: MonitorType; label: string }[] = [
    { value: "http", label: "HTTP" },
    { value: "ping", label: "Ping" },
  ];

  return (
    <form
      id={formId}
      onSubmit={(e) => {
        e.preventDefault();
        form.handleSubmit();
      }}
    >
      <Flex direction="column">
        <Flex direction="column" className="flex-1 gap-2">
          <form.Field
            name={"type"}
            listeners={{
              onChange: ({ value }) => {
                if (isValidMonitorType(value)) {
                  form.setFieldValue("config", defaultConfigs[value]);
                }
              },
            }}
            children={(field) => {
              return (
                <Field>
                  <FieldLabel>Type</FieldLabel>
                  <LMSelect
                    id={field.name}
                    name={field.name}
                    value={field.state.value}
                    onValueChange={(value) => {
                      if (isValidMonitorType(value)) {
                        field.handleChange(value as MonitorType);
                      }
                    }}
                    placeholder="Select Monitor Type"
                    items={monitorTypeSelectItems}
                    isInvalid={isFieldInvalid(field)}
                    errorMessage={getFirstError(field)}
                  />
                </Field>
              );
            }}
          />
          <form.Field
            name="name"
            listeners={{
              onChange: ({ value }) => {
                if (!useCustomSlug) {
                  form.setFieldValue("slug", slugFromString(value));
                }
              },
            }}
            children={(field) => (
              <Field>
                <FieldLabel>Name</FieldLabel>
                <LMInputField
                  name={field.name}
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  placeholder="My Monitor"
                  isInvalid={isFieldInvalid(field)}
                  errorMessage={getFirstError(field)}
                />
              </Field>
            )}
          />
          <Flex direction="row" className="gap-2 items-center">
            <FieldLabel>Use Custom Slug</FieldLabel>
            <Switch
              checked={useCustomSlug}
              onCheckedChange={onUseCustomSlugChanged}
              name="useCustomSlug"
            />
          </Flex>
          <form.Field
            name="slug"
            listeners={{
              onChange: () => {
                if (!isSlugModified) {
                  setIsSlugModified(true);
                }
              },
            }}
            children={(field) => (
              <Field>
                <FieldLabel>Slug</FieldLabel>
                <LMInputField
                  name={field.name}
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  disabled={!useCustomSlug}
                  isInvalid={isFieldInvalid(field)}
                  errorMessage={getFirstError(field)}
                />
              </Field>
            )}
          />
          <form.Field
            name="description"
            children={(field) => (
              <Field>
                <FieldLabel>Description</FieldLabel>
                <LMTextareaField
                  name={field.name}
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  isInvalid={isFieldInvalid(field)}
                  errorMessage={getFirstError(field)}
                />
              </Field>
            )}
          />
          <form.Field
            name="interval"
            children={(field) => (
              <Field>
                <FieldLabel>Interval (s)</FieldLabel>
                <LMInputField
                  name={field.name}
                  value={field.state.value.toString()}
                  onChange={(e) => field.handleChange(Number(e.target.value))}
                  isInvalid={isFieldInvalid(field)}
                  errorMessage={getFirstError(field)}
                />
              </Field>
            )}
          />
        </Flex>
        <form.Subscribe selector={(form) => form.values.type}>
          {(type) => {
            if (!type) return null;
            return <Divider direction="row" className="my-4" />;
          }}
        </form.Subscribe>

        <Flex direction="column" className="flex-1 gap-2">
          <MonitorConfigFields form={form} />
        </Flex>
      </Flex>
    </form>
  );
}
