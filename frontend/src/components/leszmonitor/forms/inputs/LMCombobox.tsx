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
    <ErrorTooltip
      isOpen={props.isInvalid ?? false}
      message={props.errorMessage ?? ""}
    >
      <Combobox
        items={props.items}
        value={props.value}
        onValueChange={props.onValueChange}
      >
        <ComboboxInput
          placeholder={props.placeholder ?? ""}
          id={props.id}
          name={props.name}
          className={cn(
            "w-full",
            props.className,
            props.isInvalid && "border-red-500 focus:ring-red-500",
          )}
          autoComplete="off"
        />
        <ComboboxContent>
          <ComboboxEmpty>{props.emptyText ?? ""}</ComboboxEmpty>
          <ComboboxList>
            {(value) => {
              return (
                <ComboboxItem key={value} value={value}>
                  {value}
                </ComboboxItem>
              );
            }}
          </ComboboxList>
        </ComboboxContent>
      </Combobox>
    </ErrorTooltip>
  );
}
