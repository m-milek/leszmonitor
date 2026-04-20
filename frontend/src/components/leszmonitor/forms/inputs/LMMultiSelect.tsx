import { ErrorTooltip } from "@/components/leszmonitor/forms/inputs/ErrorTooltip.tsx";
import {
  Combobox,
  ComboboxChip,
  ComboboxChips,
  ComboboxChipsInput,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxItem,
  ComboboxList,
  ComboboxValue,
  useComboboxAnchor,
} from "@/components/ui/combobox";
import React from "react";

export interface LMMultiSelectProps {
  value?: string[];
  onChange: (value: string[]) => void;
  options: string[];
  name: string;
  placeholder?: string;
  emptyMessage?: string;
  isInvalid?: boolean;
  errorMessage?: string;
}

export function LMMultiSelect(props: LMMultiSelectProps) {
  const anchor = useComboboxAnchor();

  const shouldDisplayPlaceholder = !props.value || props.value.length === 0;

  return (
    <ErrorTooltip isOpen={props.isInvalid} message={props.errorMessage}>
      <Combobox
        items={props.options}
        multiple
        value={props.value ?? []}
        onValueChange={props.onChange}
      >
        <ComboboxChips ref={anchor}>
          <ComboboxValue>
            {(values) => (
              <React.Fragment>
                {values.map((value: string) => (
                  <ComboboxChip key={value}>{value}</ComboboxChip>
                ))}
                <ComboboxChipsInput
                  placeholder={
                    shouldDisplayPlaceholder ? props.placeholder : ""
                  }
                />
              </React.Fragment>
            )}
          </ComboboxValue>
        </ComboboxChips>
        <ComboboxContent anchor={anchor}>
          <ComboboxEmpty>
            {props.emptyMessage ?? "No items found."}
          </ComboboxEmpty>
          <ComboboxList>
            {(item) => (
              <ComboboxItem key={item} value={item}>
                {item}
              </ComboboxItem>
            )}
          </ComboboxList>
        </ComboboxContent>
      </Combobox>
    </ErrorTooltip>
  );
}
