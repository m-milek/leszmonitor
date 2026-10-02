import type { ReactNode } from "react";
import { SearchIcon } from "lucide-react";
import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
  ComboboxTrigger,
} from "@/components/ui/combobox";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";

export interface LMFacetedFilterOption {
  value: string;
  label: string;
  render?: ReactNode;
}

export interface LMFacetedFilterProps {
  title: string;
  options: LMFacetedFilterOption[];
  value: string[];
  onChange: (value: string[]) => void;
  emptyMessage?: string;
}

export function LMFacetedFilter(props: Readonly<LMFacetedFilterProps>) {
  const selectedOptions = props.options.filter((option) =>
    props.value.includes(option.value),
  );

  return (
    <Combobox
      items={props.options}
      multiple
      value={selectedOptions}
      onValueChange={(next) => props.onChange(next.map((o) => o.value))}
      isItemEqualToValue={(a, b) => a.value === b.value}
    >
      <ComboboxTrigger
        render={<Button variant="outline" className="font-normal" />}
      >
        <SearchIcon />
        {props.title}
        {selectedOptions.length > 0 && (
          <>
            <Separator orientation="vertical" className="mx-1" />
            <Badge variant="secondary" className="font-normal">
              {selectedOptions.length} selected
            </Badge>
          </>
        )}
      </ComboboxTrigger>
      <ComboboxContent className="w-56">
        <ComboboxInput showTrigger={false} placeholder={props.title} />
        <ComboboxEmpty>{props.emptyMessage ?? "No results."}</ComboboxEmpty>
        <ComboboxList>
          {(option: LMFacetedFilterOption) => (
            <ComboboxItem key={option.value} value={option}>
              {option.render ?? option.label}
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  );
}
