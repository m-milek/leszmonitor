import { HttpMonitorConfigFields } from "@/features/monitors/forms/fields/HttpMonitorConfigFields";
import { TcpMonitorConfigFields } from "@/features/monitors/forms/fields/TcpMonitorConfigFields";
import type { MonitorFormApi } from "@/features/monitors/hooks/useMonitorForm";
import { DnsMonitorConfigFields } from "@/features/monitors/forms/fields/DnsMonitorConfigFields";
import { PushMonitorConfigFields } from "@/features/monitors/forms/fields/PushMonitorConfigFields";
import { PingMonitorConfigFields } from "@/features/monitors/forms/fields/PingMonitorConfigFields";

export function MonitorConfigFields({
  form,
}: Readonly<{ form: MonitorFormApi }>) {
  return (
    <form.Subscribe selector={(state) => state.values.type}>
      {(type) => {
        switch (type) {
          case "http":
            return <HttpMonitorConfigFields form={form} />;
          case "tcp":
            return <TcpMonitorConfigFields form={form} />;
          case "dns":
            return <DnsMonitorConfigFields form={form} />;
          case "push":
            return <PushMonitorConfigFields form={form} />;
          case "ping":
            return <PingMonitorConfigFields form={form} />;
          default:
            return null;
        }
      }}
    </form.Subscribe>
  );
}
