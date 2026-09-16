import { z } from "zod";
import { useForm } from "@tanstack/react-form";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field.tsx";
import { LMInputField } from "@/components/form/LMInputField.tsx";
import { LMColorPicker } from "@/components/form/LMColorPicker.tsx";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/form/field-state.ts";
import { Tag } from "@/features/tags/components/Tag.tsx";
import { Flex } from "@/components/common/Flex.tsx";
import {
  isValidHexColor,
  normalizeHexColor,
  randomTagColor,
} from "@/features/tags/lib/colors.ts";
import type { TagPayload } from "@/features/tags/api/tags.ts";

export interface TagFormProps {
  id?: string;
  defaultValues?: Partial<TagPayload>;
  onSubmit: (values: TagPayload) => Promise<void>;
  className?: string;
}

const formSchema = z.object({
  name: z.string().min(1, "Tag name cannot be empty"),
  description: z.string(),
  colorHex: z
    .string()
    .refine(isValidHexColor, "Color has to be a hex code, e.g. #3b82f6"),
});

export function TagForm({
  id = "tag-form",
  defaultValues,
  onSubmit,
  className,
}: Readonly<TagFormProps>) {
  const form = useForm({
    defaultValues: {
      name: defaultValues?.name ?? "",
      description: defaultValues?.description ?? "",
      colorHex: defaultValues?.colorHex ?? randomTagColor(),
    },
    validators: {
      onSubmit: formSchema,
    },
    onSubmit: async ({ value }) => {
      await onSubmit({
        name: value.name.trim(),
        description: value.description.trim(),
        colorHex: normalizeHexColor(value.colorHex),
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
      <FieldGroup className="gap-3">
        <form.Field
          name="name"
          children={(field) => (
            <Field id={field.name}>
              <FieldLabel>Name</FieldLabel>
              <LMInputField
                name={field.name}
                type="text"
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
            <Field id={field.name}>
              <FieldLabel>Description</FieldLabel>
              <LMInputField
                name={field.name}
                type="text"
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                isInvalid={isFieldInvalid(field)}
                errorMessage={getFirstError(field)}
              />
            </Field>
          )}
        />
        <form.Field
          name="colorHex"
          children={(field) => (
            <Field id={field.name}>
              <FieldLabel>Color</FieldLabel>
              <LMColorPicker
                name={field.name}
                value={field.state.value}
                onChange={(colorHex) => field.handleChange(colorHex)}
                isInvalid={isFieldInvalid(field)}
                errorMessage={getFirstError(field)}
              />
            </Field>
          )}
        />
        <form.Subscribe
          selector={(state) => [state.values.name, state.values.colorHex]}
          children={([name, colorHex]) => (
            <Flex className="justify-center pt-2">
              <Tag tag={{ name: name || "leszmonitor", colorHex }} />
            </Flex>
          )}
        />
      </FieldGroup>
    </form>
  );
}
