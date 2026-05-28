import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { Field, FieldLabel } from "@/components/ui/field.tsx";
import { LMInputField } from "@/components/leszmonitor/forms/inputs/LMInputField.tsx";
import { LMSelect } from "@/components/leszmonitor/forms/inputs/LMSelect.tsx";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/leszmonitor/forms/inputs/utils.ts";
import type { DnsRecordType } from "@/lib/types.ts";
import type { MonitorFormApi } from "@/hooks/useMonitorForm";

const dnsRecordOptions = [
  { value: "A", label: "A" },
  { value: "AAAA", label: "AAAA" },
  { value: "CNAME", label: "CNAME" },
  { value: "MX", label: "MX" },
  { value: "TXT", label: "TXT" },
  { value: "SRV", label: "SRV" },
  { value: "NS", label: "NS" },
];

export function DnsMonitorConfigFields({
  form,
}: Readonly<{ form: MonitorFormApi }>) {
  return (
    <Flex direction="column" className="gap-4 items-stretch">
      <div className="text-lg font-semibold">DNS Settings</div>

      <form.Field
        name="probeConfig.hostname"
        children={(field) => (
          <Field>
            <FieldLabel>Hostname</FieldLabel>
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
        name="probeConfig.recordType"
        children={(field) => (
          <Field>
            <FieldLabel>Record Type</FieldLabel>
            <LMSelect
              id={field.name}
              name={field.name}
              value={field.state.value}
              onValueChange={(value) =>
                field.handleChange(value as DnsRecordType)
              }
              placeholder="Select Record Type"
              items={dnsRecordOptions}
              isInvalid={isFieldInvalid(field)}
              errorMessage={getFirstError(field)}
            />
          </Field>
        )}
      />
    </Flex>
  );
}
