import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx";
import { ErrorTooltip } from "@/components/leszmonitor/forms/inputs/ErrorTooltip.tsx";
import { cn } from "@/lib/utils.ts";

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

export function LMSelect(props: LMSelectProps) {
  return (
    <Select
      value={props.value}
      onValueChange={props.onValueChange}
      autoComplete="off"
    >
      <ErrorTooltip
        isOpen={props.isInvalid ?? false}
        message={props.errorMessage ?? ""}
      >
        <SelectTrigger
          className={cn(
            props.className,
            props.isInvalid && "border-destructive focus:ring-destructive",
          )}
        >
          <SelectValue placeholder={props.placeholder} />
        </SelectTrigger>
      </ErrorTooltip>
      <SelectContent position="popper">
        {props.items?.map((option) => (
          <SelectItem key={option.value} value={option.value}>
            {option.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
