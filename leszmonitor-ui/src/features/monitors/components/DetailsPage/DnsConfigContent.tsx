import type { DnsMonitor, DnsMonitorConfig } from "@/features/monitors/types";
import { ConfigTable } from "@/features/monitors/components/DetailsPage/ConfigTable";
import {
  type FieldConfig,
  formatList,
} from "@/features/monitors/components/DetailsPage/config-fields";

export interface DnsConfigContentProps {
  monitor: DnsMonitor;
}

const dnsFields: FieldConfig<DnsMonitorConfig> = {
  hostname: { label: "Hostname" },
  recordType: { label: "Record Type" },
  dnsServer: { label: "DNS Server" },
  expectedRecordValues: {
    label: "Expected Record Values",
    render: formatList,
  },
};

export const DnsConfigContent = ({ monitor }: DnsConfigContentProps) => {
  if (!monitor.probeConfig) return null;
  return <ConfigTable data={monitor.probeConfig} fields={dnsFields} />;
};
