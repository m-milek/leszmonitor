import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "@/components/ui/combobox.tsx";
import { cn } from "@/lib/utils.ts";
import { ErrorTooltip } from "@/components/leszmonitor/forms/inputs/ErrorTooltip.tsx";

export interface LMComboboxProps {
  id: string;
  name: string;
  items: string[];
  value?: string;
  onValueChange?: (value: string | null, eventDetails: unknown) => void;
  emptyText?: string;
  placeholder?: string;
  className?: string;
  isInvalid?: boolean;
  errorMessage?: string;
}

export function LMCombobox(props: LMComboboxProps) {
  return (
    <Combobox
      items={props.items}
      value={props.value}
      onValueChange={props.onValueChange}
    >
      <ErrorTooltip
        isOpen={props.isInvalid ?? false}
        message={props.errorMessage ?? ""}
      >
        <ComboboxInput
          placeholder={props.placeholder ?? ""}
          id={props.id}
          name={props.name}
          className={cn(
            props.className,
            props.isInvalid && "border-destructive focus:ring-destructive",
          )}
          autoComplete="off"
        />
      </ErrorTooltip>
      <ComboboxContent>
        <ComboboxEmpty>{props.emptyText ?? ""}</ComboboxEmpty>
        <ComboboxList>
          {(value) => (
            <ComboboxItem key={value} value={value}>
              {value}
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  );
}
