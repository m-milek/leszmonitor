import { useForm } from "@tanstack/react-form";
import {
  isValidMonitorType,
  type MonitorType,
  newMonitorSchema,
  newMonitorSchemaDefaultValues,
} from "@/lib/types.ts";
import { Field, FieldLabel } from "@/components/ui/field.tsx";
import { LMInputField } from "@/components/leszmonitor/forms/inputs/LMInputField.tsx";
import type { ReactNode } from "react";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { LMSelect } from "@/components/leszmonitor/forms/inputs/LMSelect.tsx";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/leszmonitor/forms/inputs/utils.ts";
import { LMTextareaField } from "@/components/leszmonitor/forms/inputs/LMTextareaField.tsx";
import { Divider } from "@/components/leszmonitor/ui/Divider.tsx";

export interface MonitorFormValues {
  name: string;
  displayId: string;
  interval: string;
  projectId: string;
  type: MonitorType;
}

export interface NewMonitorFormProps {
  onSubmitMonitor: (value: MonitorFormValues | unknown) => Promise<unknown>;
  projects: Array<{ id: string; name: string }> | undefined;
  formId?: string;
  renderMonitorTypeContent?: (type: MonitorType | null) => ReactNode;
}

export function NewMonitorForm({
  onSubmitMonitor,
  projects,
  formId = "new-monitor-form",
}: NewMonitorFormProps) {
  const form = useForm({
    defaultValues: newMonitorSchemaDefaultValues,
    validators: {
      onSubmit: newMonitorSchema,
    },
    onSubmit: async ({ value }) => {
      await onSubmitMonitor(value);
    },
    onSubmitInvalid: ({ value }) => {
      console.log("Invalid form submission");
      console.log("Values:", value);
    },
  });

  const projectSelectItems = projects
    ? projects.map((project) => ({
        value: project.id,
        label: project.name,
      }))
    : [];

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
      <Flex direction="row">
        <Flex direction="column" className="flex-1 gap-2">
          <form.Field
            name="name"
            children={(field) => (
              <Field>
                <FieldLabel>Monitor Name</FieldLabel>
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
          <form.Field
            name="displayId"
            children={(field) => (
              <Field>
                <FieldLabel>Slug</FieldLabel>
                <LMInputField
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
                  value={field.state.value}
                  onChange={(e) => field.handleChange(Number(e.target.value))}
                  isInvalid={isFieldInvalid(field)}
                  errorMessage={getFirstError(field)}
                />
              </Field>
            )}
          />
          <form.Field
            name={"projectId"}
            children={(field) => {
              return (
                <Field>
                  <FieldLabel>Project</FieldLabel>
                  <LMSelect
                    id={field.name}
                    name={field.name}
                    value={field.state.value}
                    onValueChange={(value) => field.handleChange(value)}
                    placeholder="Select Project"
                    items={projectSelectItems}
                    isInvalid={isFieldInvalid(field)}
                    errorMessage={getFirstError(field)}
                  />
                </Field>
              );
            }}
          />
          <form.Field
            name={"type"}
            children={(field) => {
              return (
                <Field>
                  <FieldLabel>Monitor Type</FieldLabel>
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
        </Flex>
        <Divider direction="column" className="mx-4" />
        <Flex direction="column" className="flex-1 gap-2"></Flex>
      </Flex>
    </form>
  );
}
