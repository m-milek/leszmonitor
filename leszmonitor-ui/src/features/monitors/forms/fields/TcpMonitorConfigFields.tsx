import { Flex } from "@/components/common/Flex";
import { Field, FieldLabel } from "@/components/ui/field";
import { LMInputField } from "@/components/form/LMInputField";
import { LMSelect } from "@/components/form/LMSelect";
import { getFirstError, isFieldInvalid } from "@/components/form/field-state";
import type { TcpProtocol } from "@/features/monitors/types";
import type { MonitorFormApi } from "@/features/monitors/hooks/useMonitorForm";

const protocolItems: { value: TcpProtocol; label: string }[] = [
  { value: "tcp", label: "TCP" },
  { value: "tcp4", label: "TCP (Force IPv4)" },
  { value: "tcp6", label: "TCP (Force IPv6)" },
];

export function TcpMonitorConfigFields({
  form,
}: Readonly<{ form: MonitorFormApi }>) {
  return (
    <Flex direction="column" className="gap-4 items-stretch">
      <div className="text-lg font-semibold">TCP Monitor Settings</div>

      <form.Field name="probeConfig.host">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Host</FieldLabel>
            <LMInputField
              name={field.name}
              value={field.state.value}
              onChange={(e) => field.handleChange(e.target.value)}
              placeholder="example.com"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="probeConfig.port">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Port</FieldLabel>
            <LMInputField
              name={field.name}
              type="number"
              inputMode="numeric"
              value={field.state.value}
              onChange={(e) => field.handleChange(Number(e.target.value))}
              placeholder="443"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="probeConfig.protocol">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Protocol</FieldLabel>
            <LMSelect
              id={field.name}
              name={field.name}
              value={field.state.value}
              onValueChange={(value) =>
                field.handleChange(value as TcpProtocol)
              }
              placeholder="Select Protocol"
              items={protocolItems}
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="probeConfig.timeout">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Timeout (ms)</FieldLabel>
            <LMInputField
              name={field.name}
              type="number"
              inputMode="numeric"
              value={field.state.value}
              onChange={(e) => field.handleChange(Number(e.target.value))}
              placeholder="5000"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>

      <form.Field name="probeConfig.retryCount">
        {(field) => (
          <Field id={field.name}>
            <FieldLabel>Retry Count</FieldLabel>
            <LMInputField
              name={field.name}
              type="number"
              inputMode="numeric"
              value={field.state.value}
              onChange={(e) => field.handleChange(Number(e.target.value))}
              placeholder="3"
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      </form.Field>
    </Flex>
  );
}
