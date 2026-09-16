import { ErrorTooltip } from "@/components/form/ErrorTooltip";
import { Switch } from "@/components/ui/switch";

export interface LMSwitchProps {
  name: string;
  checked: boolean;
  onCheckedChange: (checked: boolean) => void;
  isInvalid?: boolean;
  errorMessage?: string;
}

export function LMSwitch(props: Readonly<LMSwitchProps>) {
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
