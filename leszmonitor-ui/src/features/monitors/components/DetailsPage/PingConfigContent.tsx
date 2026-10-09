import type { PingMonitor, PingMonitorConfig } from "@/features/monitors/types";
import { ConfigTable } from "@/features/monitors/components/DetailsPage/ConfigTable";
import {
  type FieldConfig,
  formatMs,
} from "@/features/monitors/components/DetailsPage/config-fields";

export interface PingConfigContentProps {
  monitor: PingMonitor;
}

const pingFields: FieldConfig<PingMonitorConfig> = {
  host: { label: "Host" },
  count: { label: "Packets" },
  timeoutMs: { label: "Timeout", render: formatMs },
};

export const PingConfigContent = ({ monitor }: PingConfigContentProps) => {
  if (!monitor.probeConfig) return null;
  return <ConfigTable data={monitor.probeConfig} fields={pingFields} />;
};
