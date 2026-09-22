import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ErrorTooltip } from "@/components/form/ErrorTooltip";

export interface LMSelectOption {
  value: string;
  label: string;
}

export interface LMSelectProps {
  id: string;
  name: string;
  value?: string;
  onValueChange?: (value: string) => void;
  placeholder?: string;
  items?: LMSelectOption[];
  className?: string;
  isInvalid?: boolean;
  errorMessage?: string;
}

export function LMSelect(props: Readonly<LMSelectProps>) {
  return (
    <Select
      value={props.value}
      onValueChange={(value) => props.onValueChange?.(value ?? "")}
    >
      <ErrorTooltip
        isOpen={props.isInvalid ?? false}
        message={props.errorMessage ?? ""}
      >
        <SelectTrigger
          id={props.id}
          className={props.className}
          aria-invalid={props.isInvalid}
        >
          <SelectValue placeholder={props.placeholder} />
        </SelectTrigger>
      </ErrorTooltip>
      <SelectContent alignItemWithTrigger={false}>
        {props.items?.map((option) => (
          <SelectItem key={option.value} value={option.value}>
            {option.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
