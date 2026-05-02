import type {
  MonitorFormApi,
} from "@/lib/formTypes.ts";
import { HttpMonitorConfigFields } from "@/components/leszmonitor/forms/monitors/HttpMonitorConfigFields.tsx";
import { PingMonitorConfigFields } from "@/components/leszmonitor/forms/monitors/PingMonitorConfigFields.tsx";

export function MonitorConfigFields({ form }: { form: MonitorFormApi }) {
  return (
    <form.Subscribe
      selector={(state) => state.values.type}
      children={(type) => {
        switch (type) {
          case "http":
            return (
              <HttpMonitorConfigFields
                form={form}
              />
            );
          case "ping":
            return (
              <PingMonitorConfigFields
                form={form}
              />
            );
          default:
            return null;
        }
      }}
    />
  );
}
