import { HttpMonitorConfigFields } from "@/features/monitors/forms/fields/HttpMonitorConfigFields.tsx";
import { TcpMonitorConfigFields } from "@/features/monitors/forms/fields/TcpMonitorConfigFields.tsx";
import type { MonitorFormApi } from "@/features/monitors/hooks/useMonitorForm.ts";
import { DnsMonitorConfigFields } from "@/features/monitors/forms/fields/DnsMonitorConfigFields.tsx";

export function MonitorConfigFields({
  form,
}: Readonly<{ form: MonitorFormApi }>) {
  return (
    <form.Subscribe
      selector={(state) => state.values.type}
      children={(type) => {
        switch (type) {
          case "http":
            return <HttpMonitorConfigFields form={form} />;
          case "tcp":
            return <TcpMonitorConfigFields form={form} />;
          case "dns":
            return <DnsMonitorConfigFields form={form} />;
          default:
            return null;
        }
      }}
    />
  );
}
