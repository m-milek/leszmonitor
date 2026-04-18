import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx";

export interface LMSelectOption {
  value: string;
  label: string;
}

export interface LMSelectProps {
  value?: string;
  onValueChange?: (value: string) => void;
  placeholder?: string;
  items?: LMSelectOption[];
}

export function LMSelect(props: LMSelectProps) {
  return (
    <Select
      value={props.value}
      onValueChange={props.onValueChange}
      autoComplete="off"
    >
      <SelectTrigger className="w-full">
        <SelectValue placeholder={props.placeholder} />
      </SelectTrigger>
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
