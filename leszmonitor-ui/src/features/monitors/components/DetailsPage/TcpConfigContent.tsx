import type { TcpMonitor, TcpMonitorConfig } from "@/features/monitors/types";
import { ConfigTable } from "@/features/monitors/components/DetailsPage/ConfigTable";
import {
  type FieldConfig,
  formatMs,
} from "@/features/monitors/components/DetailsPage/config-fields";

export interface TcpConfigContentProps {
  monitor: TcpMonitor;
}

const tcpFields: FieldConfig<TcpMonitorConfig> = {
  host: {
    label: "Hostname",
    render: (host, config) => (
      <pre>
        {host}:{config.port}
      </pre>
    ),
  },
  protocol: { label: "Protocol" },
  timeout: { label: "Timeout", render: formatMs },
};

export const TcpConfigContent = ({ monitor }: TcpConfigContentProps) => {
  if (!monitor.probeConfig) return null;
  return <ConfigTable data={monitor.probeConfig} fields={tcpFields} />;
};
