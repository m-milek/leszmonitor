import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "@/components/ui/combobox.tsx";
import { cn } from "@/lib/utils.ts";

export interface LMComboboxProps {
  id: string;
  name: string;
  items: string[];
  value?: string;
  onValueChange?: (value: string | null, eventDetails: unknown) => void;
  emptyText?: string;
  placeholder?: string;
  className?: string;
}

export function LMCombobox(props: LMComboboxProps) {
  return (
    <Combobox
      items={props.items}
      value={props.value}
      onValueChange={props.onValueChange}
    >
      <ComboboxInput
        placeholder={props.placeholder ?? ""}
        id={props.id}
        name={props.name}
        className={cn("w-full", props.className)}
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
  );
}
