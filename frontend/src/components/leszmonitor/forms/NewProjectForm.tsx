import z from "zod";
import { FieldGroup } from "@/components/ui/field.tsx";
import { useForm } from "@tanstack/react-form";
import type { ProjectInput } from "@/lib/data/projectData.ts";
import { LMInputField } from "@/components/leszmonitor/forms/LMInputField.tsx";

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
      <FieldGroup className="gap-2">
        <form.Field
          name="name"
          children={(field) => (
            <LMInputField label="Project Name" field={field} type="text" />
          )}
        />
        <form.Field
          name="displayId"
          children={(field) => (
            <LMInputField label="Display ID" field={field} type="text" />
          )}
        />
        <form.Field
          name="description"
          children={(field) => (
            <LMInputField
              label="Description (Optional)"
              field={field}
              type="textarea"
              rows={4}
            />
          )}
        />
      </FieldGroup>
    </form>
  );
}
