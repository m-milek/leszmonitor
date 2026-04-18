import { useForm } from "@tanstack/react-form";
import {
  isValidMonitorType,
  type MonitorType,
  newMonitorSchema,
  newMonitorSchemaDefaultValues,
} from "@/lib/types.ts";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx";
import { Field, FieldError, FieldLabel } from "@/components/ui/field.tsx";
import { LMInputField } from "@/components/leszmonitor/forms/LMInputField.tsx";
import type { ReactNode } from "react";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";

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
              return <LMInputField label="Display ID" field={field} />;
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
              const isInvalid =
                field.state.meta.isTouched && !field.state.meta.isValid;
              return (
                <Field>
                  <FieldLabel>Project</FieldLabel>
                  <Select
                    value={field.state.value}
                    onValueChange={(value) => field.handleChange(value)}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select Project" />
                    </SelectTrigger>
                    <SelectContent>
                      {projects?.map((project) => (
                        <SelectItem key={project.id} value={project.id}>
                          {project.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {isInvalid && <FieldError errors={field.state.meta.errors} />}
                </Field>
              );
            }}
          />
          <form.Field
            name={"type"}
            children={(field) => {
              const isInvalid =
                field.state.meta.isTouched && !field.state.meta.isValid;
              return (
                <Field>
                  <FieldLabel>Type</FieldLabel>
                  <Select
                    value={field.state.value}
                    onValueChange={(value) => {
                      if (isValidMonitorType(value)) {
                        field.handleChange(value as MonitorType);
                      }
                    }}
                  >
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Select Monitor Type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value={"http"}>HTTP</SelectItem>
                      <SelectItem value={"ping"}>Ping</SelectItem>
                    </SelectContent>
                  </Select>
                  {isInvalid && <FieldError errors={field.state.meta.errors} />}
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


