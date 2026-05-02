import { type FormValidateOrFn, useForm } from "@tanstack/react-form";
import { useState } from "react";
import { slugFromString } from "@/lib/slugFromString.ts";
import type { MonitorFormValues } from "@/lib/types.ts";
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
import { Switch } from "@/components/ui/switch.tsx";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createMonitor } from "@/lib/data/monitorData.ts";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { toast } from "sonner";

export interface NewMonitorFormProps {
  projectSlug: string;
  formId?: string;
}

export interface MonitorFormProps {
  projectSlug: string;
  formId?: string;
  defaultValues?: Partial<MonitorFormValues>;
  onSubmit: (value: MonitorFormValues) => Promise<void>;
  resetOnSuccess?: boolean;
}

const buildMonitorDefaults = (
  projectSlug: string,
  defaultValues?: Partial<MonitorFormValues>,
): MonitorFormValues => {
  const baseValues = {
    ...newMonitorSchemaDefaultValues,
    projectSlug: projectSlug,
  };
  const type = defaultValues?.type;

  if (!type) {
    return {
      ...baseValues,
      ...defaultValues,
    } as MonitorFormValues;
  }

  return {
    ...baseValues,
    ...defaultValues,
    type,
    projectSlug: projectSlug,
    config: {
      ...defaultConfigs[type],
      ...(defaultValues?.config ?? {}),
    },
  } as MonitorFormValues;
};

export function MonitorForm({
  projectSlug,
  formId = "monitor-form",
  defaultValues,
  onSubmit,
  resetOnSuccess = false,
}: MonitorFormProps) {
  const mergedDefaults = buildMonitorDefaults(projectSlug, defaultValues);

  const form = useForm({
    defaultValues: mergedDefaults,
    validators: {
      onSubmit:
        newMonitorSchema as unknown as FormValidateOrFn<MonitorFormValues>,
    },
    onSubmit: async ({ value }) => {
      await onSubmit(value);
      if (resetOnSuccess) {
        form.reset();
      }
    },
    onSubmitInvalid: ({ value }) => {
      console.log("Invalid form submission");
      console.log("Values:", value);
      toast.error("Please fix the errors in the form before submitting.");
    },
  });

  const [useCustomSlug, setUseCustomSlug] = useState(() => {
    const name = mergedDefaults.name ?? "";
    const slug = mergedDefaults.slug ?? "";
    return slug.length > 0 && slug !== slugFromString(name);
  });

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

export function NewMonitorForm({
  projectSlug,
  formId = "new-monitor-form",
}: NewMonitorFormProps) {
  const queryClient = useQueryClient();
  const createMonitorMutation = useMutation({
    mutationFn: (monitor: MonitorFormValues) => {
      console.log("Creating monitor with values:", monitor);
      return createMonitor(monitor);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.MONITORS] });
    },
  });

  return (
    <MonitorForm
      projectSlug={projectSlug}
      formId={formId}
      onSubmit={(value) => createMonitorMutation.mutateAsync(value)}
      resetOnSuccess
    />
  );
}
