import { RotateCcw } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Field,
  FieldContent,
  FieldDescription,
  FieldLabel,
} from "@/components/ui/field";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { GlobalParameterInput } from "@/features/instance/components/GlobalParameterInput";
import { GlobalParameterSourceBadge } from "@/features/instance/components/GlobalParameterSourceBadge";
import type {
  GlobalParameter,
  GlobalParameterValue,
} from "@/features/instance/global-parameters-api";

export interface GlobalParameterRowProps {
  id: string;
  parameter: GlobalParameter;
  value: GlobalParameterValue | null;
  onChange: (value: GlobalParameterValue | null) => void;
}

export function GlobalParameterRow({
  id,
  parameter,
  value,
  onChange,
}: Readonly<GlobalParameterRowProps>) {
  const isLocked = parameter.source === "config";
  const canReset = parameter.source === "ui" && value !== null;

  return (
    <Field id={id} orientation="horizontal">
      <FieldContent>
        <FieldLabel>{parameter.displayName}</FieldLabel>
        <FieldDescription>{parameter.description}</FieldDescription>
      </FieldContent>
      <div className="flex w-36 shrink-0 justify-center">
        <GlobalParameterSourceBadge source={parameter.source} />
      </div>
      <div className="flex w-80 shrink-0 items-center justify-end gap-2">
        {canReset && (
          <Tooltip>
            <TooltipTrigger
              render={
                <Button
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  aria-label="Reset to default"
                  onClick={() => onChange(null)}
                />
              }
            >
              <RotateCcw />
            </TooltipTrigger>
            <TooltipContent>Reset to default</TooltipContent>
          </Tooltip>
        )}
        <GlobalParameterInput
          id={id}
          parameter={parameter}
          value={value ?? parameter.defaultValue}
          onChange={onChange}
          disabled={isLocked}
        />
      </div>
    </Field>
  );
}
