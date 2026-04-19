import type {
  MonitorFormApi,
  HttpMonitorFormApi,
  PingMonitorFormApi,
} from "@/components/leszmonitor/forms/NewMonitorForm.tsx";
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
                form={form as unknown as HttpMonitorFormApi}
              />
            );
          case "ping":
            return (
              <PingMonitorConfigFields
                form={form as unknown as PingMonitorFormApi}
              />
            );
          default:
            return null;
        }
      }}
    />
  );
}
