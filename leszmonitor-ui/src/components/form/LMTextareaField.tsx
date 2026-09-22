import { Textarea } from "@/components/ui/textarea";
import type { ComponentProps } from "react";
import { ErrorTooltip } from "@/components/form/ErrorTooltip";

type TextareaProps = ComponentProps<typeof Textarea>;

interface LMTextareaFieldProps {
  placeholder?: TextareaProps["placeholder"];
  rows?: TextareaProps["rows"];
  isInvalid?: boolean;
  errorMessage?: string;
  name: string;
  value?: string | number;
  onChange: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
}

export const LMTextareaField = (props: LMTextareaFieldProps) => {
  return (
    <ErrorTooltip
      isOpen={props.isInvalid ?? false}
      message={props.errorMessage ?? ""}
    >
      <Textarea
        id={props.name}
        name={props.name}
        value={props.value}
        onChange={props.onChange}
        placeholder={props.placeholder}
        rows={props.rows}
        aria-invalid={props.isInvalid}
      />
    </ErrorTooltip>
  );
};
