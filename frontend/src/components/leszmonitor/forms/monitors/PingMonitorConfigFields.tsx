import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { Field, FieldLabel } from "@/components/ui/field.tsx";
import { LMInputField } from "@/components/leszmonitor/forms/inputs/LMInputField.tsx";
import { LMSelect } from "@/components/leszmonitor/forms/inputs/LMSelect.tsx";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/leszmonitor/forms/inputs/utils.ts";
import type { MonitorFormApi } from "@/lib/formTypes.ts";

const protocolItems = [
  { value: "tcp", label: "TCP" },
  { value: "udp", label: "UDP" },
  { value: "tcp4", label: "TCP (Force IPv4)" },
  { value: "tcp6", label: "TCP (Force IPv6)" },
  { value: "udp4", label: "UDP (Force IPv4)" },
  { value: "udp6", label: "UDP (Force IPv6)" },
];

export function PingMonitorConfigFields({
  form,
}: {
  form: MonitorFormApi;
}) {
  return (
    <Flex direction="column" className="gap-4 items-stretch">
      <div className="text-lg font-semibold">Ping Settings</div>

      <form.Field
        name="config.host"
        children={(field) => (
          <Field>
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
      />

      <form.Field
        name="config.port"
        children={(field) => (
          <Field>
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
      />

      <form.Field
        name="config.protocol"
        children={(field) => (
          <Field>
            <FieldLabel>Protocol</FieldLabel>
            <LMSelect
              id={field.name}
              name={field.name}
              value={field.state.value}
              onValueChange={(value) =>
                field.handleChange(
                  value as "tcp" | "udp" | "tcp4" | "tcp6" | "udp4" | "udp6",
                )
              }
              placeholder="Select Protocol"
              items={protocolItems}
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      />

      <form.Field
        name="config.timeout"
        children={(field) => (
          <Field>
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
      />

      <form.Field
        name="config.retryCount"
        children={(field) => (
          <Field>
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
      />
    </Flex>
  );
}
