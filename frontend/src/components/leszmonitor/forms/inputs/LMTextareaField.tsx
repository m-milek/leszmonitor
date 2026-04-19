import { Field, FieldLabel } from "@/components/ui/field.tsx";
import { Textarea } from "@/components/ui/textarea.tsx";
import type { ComponentProps } from "react";
import { cn } from "@/lib/utils.ts";
import { ErrorTooltip } from "@/components/leszmonitor/forms/inputs/ErrorTooltip.tsx";

type TextareaProps = ComponentProps<typeof Textarea>;

interface LMTextareaFieldProps {
  label: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  field: any;
  placeholder?: TextareaProps["placeholder"];
  rows?: TextareaProps["rows"];
  parseValue?: (rawValue: string) => string | number;
}

const getFirstError = (errors?: Array<{ message?: string } | undefined>) => {
  return errors?.[0]?.message ?? "";
};

export const LMTextareaField = ({
  label,
  field,
  placeholder,
  rows,
  parseValue,
}: LMTextareaFieldProps) => {
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;
  const errorMessage = getFirstError(field.state.meta.errors);

  const handleChange = (rawValue: string) => {
    if (parseValue) {
      field.handleChange(parseValue(rawValue));
      return;
    }

    const nextValue =
      typeof field.state.value === "number" ? Number(rawValue) : rawValue;
    field.handleChange(nextValue);
  };

  return (
    <Field>
      <FieldLabel htmlFor={field.name}>{label}</FieldLabel>
      <ErrorTooltip isOpen={isInvalid} message={errorMessage}>
        <Textarea
          id={field.name}
          name={field.name}
          value={field.state.value}
          onChange={(e) => handleChange(e.target.value)}
          autoComplete="off"
          placeholder={placeholder}
          rows={rows}
          className={cn(isInvalid && "border-red-500 focus:ring-red-500")}
        />
      </ErrorTooltip>
    </Field>
  );
};
