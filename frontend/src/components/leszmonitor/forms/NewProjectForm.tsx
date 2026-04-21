import z from "zod";
import { Field, FieldLabel } from "@/components/ui/field.tsx";
import { useState } from "react";
import { isSlugValid, slugFromString } from "@/lib/slugFromString.ts";
import { useForm } from "@tanstack/react-form";
import type { ProjectInput } from "@/lib/data/projectData.ts";
import { LMInputField } from "@/components/leszmonitor/forms/inputs/LMInputField.tsx";
import { LMTextareaField } from "@/components/leszmonitor/forms/inputs/LMTextareaField.tsx";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/leszmonitor/forms/inputs/utils.ts";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { Switch } from "@/components/ui/switch.tsx";

const projectFormSchema = z.object({
  name: z.string().min(1, "Project name is required"),
  slug: z
    .string()
    .min(1, "Slug is required")
    .refine(
      isSlugValid,
      "Invalid slug format. Must be lowercase, alphanumeric, and can include hyphens.",
    ),
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
      slug: "",
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

  const [useCustomSlug, setUseCustomSlug] = useState(false);

  const onUseCustomSlugChanged = (checked: boolean) => {
    setUseCustomSlug(checked);
    if (!checked) {
      const name = form.state.values.name;
      form.setFieldValue("slug", slugFromString(name));
    }
  };

  return (
    <form
      id={formId}
      onSubmit={(e) => {
        e.preventDefault();
        form.handleSubmit();
      }}
    >
      <Flex direction="column" className="gap-2">
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
            <Field>
              <FieldLabel>Description (Optional)</FieldLabel>
              <LMTextareaField
                name={field.name}
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                rows={4}
                isInvalid={isFieldInvalid(field)}
                errorMessage={getFirstError(field)}
              />
            </Field>
          )}
        />
      </Flex>
    </form>
  );
}
