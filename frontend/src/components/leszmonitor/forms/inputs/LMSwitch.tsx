import { ErrorTooltip } from "@/components/leszmonitor/forms/inputs/ErrorTooltip.tsx";
import { Switch } from "@/components/ui/switch.tsx";

export interface LMSwitchProps {
  name: string;
  checked: boolean;
  onCheckedChange: (checked: boolean) => void;
  isInvalid?: boolean;
  errorMessage?: string;
}

export function LMSwitch(props: LMSwitchProps) {
  return (
    <ErrorTooltip isOpen={props.isInvalid} message={props.errorMessage}>
      <Switch
        id={props.name}
        name={props.name}
        checked={props.checked}
        onCheckedChange={props.onCheckedChange}
      />
    </ErrorTooltip>
  );
}
