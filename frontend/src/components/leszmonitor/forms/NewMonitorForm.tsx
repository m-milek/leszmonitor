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
  renderMonitorTypeContent,
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

  const selectedMonitorType = form.state.values.type;

  const isValidType =
    selectedMonitorType && isValidMonitorType(selectedMonitorType);

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
    <Flex direction="vertical" gap="1rem" className="w-full" align="stretch">
      <form
        id={formId}
        onSubmit={(e) => {
          e.preventDefault();
          form.handleSubmit();
        }}
      >
        <Flex direction="vertical" gap="1rem" align="stretch">
          <form.Field
            name="name"
            children={(field) => {
              return <LMInputField label="Name" field={field} />;
            }}
          />
          <form.Field
            name="displayId"
            children={(field) => {
              return <LMInputField label="Slug" field={field} />;
            }}
          />
          <form.Field
            name="interval"
            children={(field) => {
              return <LMInputField label="Interval (s)" field={field} />;
            }}
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
      </form>
      {isValidType &&
        renderMonitorTypeContent &&
        renderMonitorTypeContent(selectedMonitorType)}
    </Flex>
  );
}
