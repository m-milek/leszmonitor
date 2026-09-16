import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { slugFromString } from "@/lib/slugFromString";
import {
  isValidMonitorType,
  type MonitorType,
} from "@/features/monitors/types";
import {
  type MonitorFormValues,
  defaultConfigs,
} from "@/features/monitors/schema";
import { buildMonitorDefaults } from "@/features/monitors/forms/monitor-form-defaults";
import { Field, FieldLabel, FieldTitle } from "@/components/ui/field";
import { LMInputField } from "@/components/form/LMInputField";
import { LMSelect } from "@/components/form/LMSelect";
import { LMTextareaField } from "@/components/form/LMTextareaField";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/form/field-state";
import { Flex } from "@/components/common/Flex";
import { Divider } from "@/components/common/Divider";
import { Switch } from "@/components/ui/switch";
import { MonitorConfigFields } from "@/features/monitors/forms/fields/MonitorConfigFields";
import { useMonitorForm } from "@/features/monitors/hooks/useMonitorForm";
import { LMTagSelect } from "@/features/tags/components/LMTagSelect";
import { getAllTags } from "@/features/tags/tags-api";
import { QUERY_KEYS } from "@/lib/consts";

export interface MonitorFormProps {
  formId?: string;
  defaultValues?: Partial<MonitorFormValues>;
  onSubmit: (value: MonitorFormValues) => Promise<void>;
  resetOnSuccess?: boolean;
}

export function MonitorForm({
  formId = "monitor-form",
  defaultValues,
  onSubmit,
  resetOnSuccess = false,
}: Readonly<MonitorFormProps>) {
  const mergedDefaults = buildMonitorDefaults(defaultValues);

  const { data: tags = [] } = useQuery({
    queryKey: [QUERY_KEYS.TAGS],
    queryFn: () => getAllTags(),
  });

  const form = useMonitorForm({
    defaultValues: mergedDefaults,
    onSubmit,
    onReset: resetOnSuccess ? () => form.reset() : undefined,
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
    { value: "tcp", label: "TCP" },
    { value: "dns", label: "DNS" },
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
                  form.setFieldValue("probeConfig", defaultConfigs[value]);
                }
              },
            }}
            children={(field) => {
              return (
                <Field id={field.name}>
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
              <Field id={field.name}>
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
            <FieldTitle>Use Custom Slug</FieldTitle>
            <Switch
              checked={useCustomSlug}
              onCheckedChange={onUseCustomSlugChanged}
              name="useCustomSlug"
            />
          </Flex>
          <form.Field
            name="slug"
            children={(field) => (
              <Field id={field.name}>
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
              <Field id={field.name}>
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
            name="tagIds"
            children={(field) => (
              <Field id={field.name}>
                <FieldLabel>Tags</FieldLabel>
                <LMTagSelect
                  id={field.name}
                  name={field.name}
                  tags={tags}
                  value={field.state.value}
                  onChange={(tagIds) => field.handleChange(tagIds)}
                  placeholder="Add tag"
                  emptyMessage="No tags found."
                  isInvalid={isFieldInvalid(field)}
                  errorMessage={getFirstError(field)}
                />
              </Field>
            )}
          />
          <form.Field
            name="interval"
            children={(field) => (
              <Field id={field.name}>
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
