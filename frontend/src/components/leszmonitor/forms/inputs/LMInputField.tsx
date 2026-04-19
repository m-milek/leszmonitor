import { Field, FieldLabel } from "@/components/ui/field.tsx";
import { Input } from "@/components/ui/input.tsx";
import type { ComponentProps } from "react";
import { cn } from "@/lib/utils.ts";
import { ErrorTooltip } from "@/components/leszmonitor/forms/inputs/ErrorTooltip.tsx";

type InputProps = ComponentProps<typeof Input>;

interface LMInputFieldProps {
  label: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  field: any;
  type?: InputProps["type"];
  autoComplete?: InputProps["autoComplete"];
  placeholder?: InputProps["placeholder"];
  inputMode?: InputProps["inputMode"];
  parseValue?: (rawValue: string) => string | number;
}

const getFirstError = (errors?: Array<{ message?: string } | undefined>) => {
  return errors?.[0]?.message ?? "";
};

export const LMInputField = ({
  label,
  field,
  type,
  autoComplete = "off",
  placeholder,
  inputMode,
  parseValue,
}: LMInputFieldProps) => {
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
        <Input
          id={field.name}
          name={field.name}
          type={type}
          value={field.state.value}
          onChange={(e) => handleChange(e.target.value)}
          autoComplete={autoComplete}
          placeholder={placeholder}
          inputMode={inputMode}
          className={cn(
            isInvalid && "border-destructive focus:ring-destructive",
          )}
        />
      </ErrorTooltip>
    </Field>
  );
};
