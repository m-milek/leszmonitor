import { Input } from "@/components/ui/input.tsx";
import type { ComponentProps } from "react";
import { cn } from "@/lib/utils.ts";
import { ErrorTooltip } from "@/components/leszmonitor/forms/inputs/ErrorTooltip.tsx";

type InputProps = ComponentProps<typeof Input>;

interface LMInputFieldProps {
  type?: InputProps["type"];
  autoComplete?: InputProps["autoComplete"];
  placeholder?: InputProps["placeholder"];
  inputMode?: InputProps["inputMode"];
  isInvalid?: boolean;
  errorMessage?: string;
  name: string;
  value?: string | number;
  onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  disabled?: boolean;
}

export const LMInputField = (props: LMInputFieldProps) => {
  return (
    <ErrorTooltip isOpen={props.isInvalid} message={props.errorMessage}>
      <Input
        id={props.name}
        name={props.name}
        type={props.type}
        value={props.value}
        onChange={props.onChange}
        autoComplete={props.autoComplete}
        placeholder={props.placeholder}
        inputMode={props.inputMode}
        disabled={props.disabled}
        className={cn(
          props.isInvalid && "border-destructive focus:ring-destructive",
        )}
      />
    </ErrorTooltip>
  );
};
