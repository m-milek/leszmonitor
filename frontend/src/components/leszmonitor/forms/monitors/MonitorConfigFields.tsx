import type { MonitorFormApi } from "@/lib/formTypes.ts";
import { HttpMonitorConfigFields } from "@/components/leszmonitor/forms/monitors/HttpMonitorConfigFields.tsx";
import { TcpMonitorConfigFields } from "@/components/leszmonitor/forms/monitors/TcpMonitorConfigFields.tsx";

export function MonitorConfigFields({ form }: { form: MonitorFormApi }) {
  return (
    <form.Subscribe
      selector={(state) => state.values.type}
      children={(type) => {
        switch (type) {
          case "http":
            return <HttpMonitorConfigFields form={form} />;
          case "tcp":
            return <TcpMonitorConfigFields form={form} />;
          default:
            return null;
        }
      }}
    />
  );
}
