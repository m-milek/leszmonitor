import { useState } from "react";
import { z } from "zod";
import { useForm } from "@tanstack/react-form";
import { useQuery } from "@tanstack/react-query";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { LMInputField } from "@/components/form/LMInputField";
import { LMTextareaField } from "@/components/form/LMTextareaField";
import { LMSwitch } from "@/components/form/LMSwitch";
import { Flex } from "@/components/common/Flex";
import { getFirstError, isFieldInvalid } from "@/components/form/field-state";
import { LMTagSelect } from "@/features/tags/components/LMTagSelect";
import { TagsApi } from "@/features/tags/tags-api";
import type { StatusPagePayload } from "@/features/status-pages/status-pages-api";
import { isSlugValid, slugFromString } from "@/lib/slugFromString";
import { QUERY_KEYS } from "@/lib/consts";

export interface StatusPageFormProps {
  id?: string;
  defaultValues?: Partial<StatusPagePayload>;
  onSubmit: (values: StatusPagePayload) => Promise<void>;
  className?: string;
}

const formSchema = z.object({
  title: z.string().min(1, "Title is required"),
  slug: z
    .string()
    .min(1, "Slug is required")
    .refine(
      isSlugValid,
      "Invalid slug format. Must be lowercase, alphanumeric, and can include hyphens.",
    ),
  description: z.string(),
  groupMonitorsByTag: z.boolean(),
  tagIds: z.array(z.string()).min(1, "Select at least one tag"),
});

export function StatusPageForm({
  id = "status-page-form",
  defaultValues,
  onSubmit,
  className,
}: Readonly<StatusPageFormProps>) {
  const { data: tags = [] } = useQuery({
    queryKey: [QUERY_KEYS.TAGS],
    queryFn: () => TagsApi.getAll(),
  });

  const [isSlugEdited, setIsSlugEdited] = useState(false);

  const form = useForm({
    defaultValues: {
      title: defaultValues?.title ?? "",
      slug: defaultValues?.slug ?? "",
      description: defaultValues?.description ?? "",
      groupMonitorsByTag: defaultValues?.groupMonitorsByTag ?? false,
      tagIds: defaultValues?.tagIds ?? [],
    },
    validators: {
      onSubmit: formSchema,
    },
    onSubmit: async ({ value }) => {
      await onSubmit({
        title: value.title.trim(),
        slug: value.slug.trim(),
        description: value.description.trim(),
        groupMonitorsByTag: value.groupMonitorsByTag,
        tagIds: value.tagIds,
      });
    },
  });

  return (
    <form
      id={id}
      className={className}
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}
    >
      <FieldGroup>
        <form.Field
          name="title"
          listeners={{
            onChange: ({ value }) => {
              if (!isSlugEdited) {
                form.setFieldValue("slug", slugFromString(value));
              }
            },
          }}
        >
          {(field) => (
            <Field id={field.name}>
              <FieldLabel>Title</FieldLabel>
              <LMInputField
                name={field.name}
                type="text"
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                placeholder="Production"
                isInvalid={isFieldInvalid(field)}
                errorMessage={getFirstError(field)}
              />
            </Field>
          )}
        </form.Field>
        <form.Field name="slug">
          {(field) => (
            <Field id={field.name}>
              <FieldLabel>Slug</FieldLabel>
              <LMInputField
                name={field.name}
                type="text"
                value={field.state.value}
                onChange={(e) => {
                  setIsSlugEdited(true);
                  field.handleChange(e.target.value);
                }}
                isInvalid={isFieldInvalid(field)}
                errorMessage={getFirstError(field)}
              />
            </Field>
          )}
        </form.Field>
        <form.Field name="description">
          {(field) => (
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
        </form.Field>
        <form.Field name="tagIds">
          {(field) => (
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
        </form.Field>
        <form.Field name="groupMonitorsByTag">
          {(field) => (
            <Field id={field.name}>
              <Flex direction="row" className="justify-between">
                <FieldLabel>Group Monitors by Tag</FieldLabel>
                <LMSwitch
                  name={field.name}
                  checked={field.state.value}
                  onCheckedChange={(checked) => field.handleChange(checked)}
                />
              </Flex>
            </Field>
          )}
        </form.Field>
      </FieldGroup>
    </form>
  );
}
