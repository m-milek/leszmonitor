import React from "react";
import { XIcon } from "lucide-react";
import { ErrorTooltip } from "@/components/form/ErrorTooltip.tsx";
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
import { Tag } from "@/features/tags/components/Tag.tsx";
import { tagChipStyle } from "@/features/tags/lib/colors.ts";
import type { Tag as TagModel } from "@/features/tags/types.ts";

interface TagOption {
  value: string;
  label: string;
  tag: TagModel;
}

export interface LMTagSelectProps {
  id?: string;
  name: string;
  tags: TagModel[];
  value?: string[];
  onChange: (tagIds: string[]) => void;
  placeholder?: string;
  emptyMessage?: string;
  isInvalid?: boolean;
  errorMessage?: string;
}

export function LMTagSelect(props: Readonly<LMTagSelectProps>) {
  const anchor = useComboboxAnchor();

  const options: TagOption[] = props.tags.map((tag) => ({
    value: tag.id,
    label: tag.name,
    tag,
  }));

  const selectedIds = props.value ?? [];
  const selectedOptions = options.filter((option) =>
    selectedIds.includes(option.value),
  );

  const shouldDisplayPlaceholder = selectedOptions.length === 0;

  const removeTag = (tagId: string) =>
    props.onChange(selectedIds.filter((id) => id !== tagId));

  return (
    <ErrorTooltip isOpen={props.isInvalid} message={props.errorMessage}>
      <Combobox
        items={options}
        multiple
        value={selectedOptions}
        onValueChange={(next) => props.onChange(next.map((o) => o.value))}
        isItemEqualToValue={(a, b) => a.value === b.value}
      >
        <ComboboxChips ref={anchor}>
          <ComboboxValue>
            {(values: TagOption[]) => (
              <React.Fragment>
                {values.map((option) => (
                  <ComboboxChip
                    key={option.value}
                    style={tagChipStyle(option.tag.colorHex)}
                    className="rounded-full border px-2 text-xs"
                    showRemove={false}
                  >
                    {option.label}
                    <button
                      type="button"
                      aria-label={`Remove tag ${option.label}`}
                      className="cursor-pointer opacity-60 hover:opacity-100"
                      onClick={(e) => {
                        e.stopPropagation();
                        removeTag(option.value);
                      }}
                    >
                      <XIcon className="size-3" />
                    </button>
                  </ComboboxChip>
                ))}
                <ComboboxChipsInput
                  id={props.id}
                  name={props.name}
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
            {props.emptyMessage ?? "No tags found."}
          </ComboboxEmpty>
          <ComboboxList>
            {(option: TagOption) => (
              <ComboboxItem key={option.value} value={option}>
                <Tag tag={option.tag} className="px-2 py-0.5 text-xs" />
              </ComboboxItem>
            )}
          </ComboboxList>
        </ComboboxContent>
      </Combobox>
    </ErrorTooltip>
  );
}
