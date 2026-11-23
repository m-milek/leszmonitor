import React from "react";
import { ErrorTooltip } from "@/components/form/ErrorTooltip";
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
import {
  notificationChannelKindLabels,
  type NotificationChannel,
} from "@/features/notifications/types";

interface ChannelOption {
  value: string;
  label: string;
  channel: NotificationChannel;
}

export interface LMNotificationChannelSelectProps {
  id?: string;
  name: string;
  channels: NotificationChannel[];
  value?: string[];
  onChange: (channelIds: string[]) => void;
  placeholder?: string;
  emptyMessage?: string;
  isInvalid?: boolean;
  errorMessage?: string;
}

export function LMNotificationChannelSelect(
  props: Readonly<LMNotificationChannelSelectProps>,
) {
  const anchor = useComboboxAnchor();

  const options: ChannelOption[] = props.channels.map((channel) => ({
    value: channel.id,
    label: channel.name,
    channel,
  }));

  const selectedIds = props.value ?? [];
  const selectedOptions = options.filter((option) =>
    selectedIds.includes(option.value),
  );

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
            {(values: ChannelOption[]) => (
              <React.Fragment>
                {values.map((option) => (
                  <ComboboxChip key={option.value}>{option.label}</ComboboxChip>
                ))}
                <ComboboxChipsInput
                  id={props.id}
                  name={props.name}
                  placeholder={
                    selectedOptions.length === 0 ? props.placeholder : ""
                  }
                />
              </React.Fragment>
            )}
          </ComboboxValue>
        </ComboboxChips>
        <ComboboxContent anchor={anchor}>
          <ComboboxEmpty>
            {props.emptyMessage ?? "No notification channels found."}
          </ComboboxEmpty>
          <ComboboxList>
            {(option: ChannelOption) => (
              <ComboboxItem key={option.value} value={option}>
                <span>{option.label}</span>
                <span className="text-muted-foreground">
                  {notificationChannelKindLabels[option.channel.kind]}
                  {!option.channel.enabled && " · disabled"}
                </span>
              </ComboboxItem>
            )}
          </ComboboxList>
        </ComboboxContent>
      </Combobox>
    </ErrorTooltip>
  );
}
