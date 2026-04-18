import { Field, FieldError, FieldLabel } from "@/components/ui/field.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Textarea } from "@/components/ui/textarea.tsx";
import type { ComponentProps } from "react";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";

type InputProps = ComponentProps<typeof Input>;
type TextareaProps = ComponentProps<typeof Textarea>;

interface LMInputFieldProps {
  label: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  field: any;
  type?: InputProps["type"] | "textarea";
  autoComplete?: InputProps["autoComplete"];
  placeholder?: InputProps["placeholder"];
  inputMode?: InputProps["inputMode"];
  rows?: TextareaProps["rows"];
  parseValue?: (rawValue: string) => string | number;
}

export const LMInputField = ({
  label,
  field,
  type,
  autoComplete = "off",
  placeholder,
  inputMode,
  rows,
  parseValue,
}: LMInputFieldProps) => {
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;
  const isTextarea = type === "textarea";

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
      <Flex direction="horizontal" gap="0.5rem" align="center" wrap="nowrap">
        {isTextarea ? (
          <Textarea
            id={field.name}
            name={field.name}
            value={field.state.value}
            onChange={(e) => handleChange(e.target.value)}
            autoComplete={autoComplete}
            placeholder={placeholder}
            rows={rows}
            className="max-w-1/2"
          />
        ) : (
          <Input
            id={field.name}
            name={field.name}
            type={type}
            value={field.state.value}
            onChange={(e) => handleChange(e.target.value)}
            autoComplete={autoComplete}
            placeholder={placeholder}
            inputMode={inputMode}
            className="max-w-1/2"
          />
        )}
        {isInvalid && <FieldError errors={field.state.meta.errors} />}
      </Flex>
    </Field>
  );
};
