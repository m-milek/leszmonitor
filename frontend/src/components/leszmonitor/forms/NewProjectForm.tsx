import z from "zod";
import { Field, FieldLabel } from "@/components/ui/field.tsx";
import { useForm } from "@tanstack/react-form";
import type { ProjectInput } from "@/lib/data/projectData.ts";
import { LMInputField } from "@/components/leszmonitor/forms/inputs/LMInputField.tsx";
import { LMTextareaField } from "@/components/leszmonitor/forms/inputs/LMTextareaField.tsx";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/leszmonitor/forms/inputs/utils.ts";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";

const projectFormSchema = z.object({
  name: z.string().min(1, "Project name is required"),
  displayId: z.string().min(1, "Display ID is required"),
  description: z.string(),
});

export interface NewProjectFormProps {
  onSubmitProject: (value: ProjectInput) => Promise<unknown>;
  formId?: string;
}

export function NewProjectForm({
  onSubmitProject,
  formId = "project-form",
}: NewProjectFormProps) {
  const form = useForm({
    defaultValues: {
      name: "",
      displayId: "",
      description: "",
    },
    validators: {
      onSubmit: projectFormSchema,
    },
    onSubmit: async ({ value }) => {
      await onSubmitProject(value);
      form.reset();
    },
  });

  return (
    <form
      id={formId}
      onSubmit={(e) => {
        e.preventDefault();
        form.handleSubmit();
      }}
    >
      <Flex direction="vertical" gap="0.5rem">
        <form.Field
          name="name"
          children={(field) => (
            <Field>
              <FieldLabel>Project Name</FieldLabel>
              <LMInputField
                name={field.name}
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                type="text"
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
              <FieldLabel>Display ID</FieldLabel>
              <LMInputField
                name={field.name}
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                type="text"
                isInvalid={isFieldInvalid(field)}
                errorMessage={getFirstError(field)}
              />
            </Field>
          )}
        />
        <form.Field
          name="description"
          children={(field) => (
            <LMTextareaField
              label="Description (Optional)"
              field={field}
              rows={4}
            />
          )}
        />
      </Flex>
    </form>
  );
}
