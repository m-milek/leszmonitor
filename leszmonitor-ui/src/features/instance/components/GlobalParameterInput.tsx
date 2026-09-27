import { LMInputField } from "@/components/form/LMInputField";
import { LMSelect } from "@/components/form/LMSelect";
import { LMSwitch } from "@/components/form/LMSwitch";
import type {
  GlobalParameter,
  GlobalParameterValue,
} from "@/features/instance/global-parameters-api";

export interface GlobalParameterInputProps {
  id: string;
  parameter: GlobalParameter;
  value: GlobalParameterValue;
  onChange: (value: GlobalParameterValue) => void;
  disabled?: boolean;
}

export function GlobalParameterInput({
  id,
  parameter,
  value,
  onChange,
  disabled,
}: Readonly<GlobalParameterInputProps>) {
  switch (parameter.type) {
    case "bool":
      return (
        <LMSwitch
          name={id}
          checked={value as boolean}
          onCheckedChange={onChange}
          disabled={disabled}
        />
      );
    case "int":
      return (
        <LMInputField
          name={id}
          type="number"
          value={value as number}
          onChange={(e) => onChange(e.target.valueAsNumber)}
          disabled={disabled}
        />
      );
    case "string":
      return (
        <LMInputField
          name={id}
          type="text"
          value={value as string}
          onChange={(e) => onChange(e.target.value)}
          disabled={disabled}
        />
      );
    case "enum":
      return (
        <LMSelect
          id={id}
          name={id}
          value={value as string}
          onValueChange={onChange}
          items={parameter.options?.map((option) => ({
            value: option,
            label: option,
          }))}
          disabled={disabled}
        />
      );
  }
}
