import { Flex } from "@/components/common/Flex";
import { StatusDot } from "@/components/common/StatusDot";
import { BatteryChart } from "@/features/monitors/components/charts/BatteryChart/BatteryChart";
import { monitorStatusToStatusDot } from "@/features/monitors/status";
import type { PublicMonitor } from "@/features/status-pages/types";

export interface PublicMonitorRowProps {
  monitor: PublicMonitor;
}

export const PublicMonitorRow = ({ monitor }: PublicMonitorRowProps) => {
  return (
    <Flex direction="column" className="gap-2 py-3">
      <Flex direction="row" className="items-center justify-between gap-4">
        <Flex direction="row" className="items-center gap-2">
          <StatusDot status={monitorStatusToStatusDot(monitor.status)} />
          <span className="font-medium">{monitor.name}</span>
        </Flex>
        <span className="text-muted-foreground text-sm">
          {monitor.uptimePercent === null
            ? "No data"
            : `${monitor.uptimePercent.toFixed(2)}% uptime`}
        </span>
      </Flex>
      <BatteryChart
        monitorResults={monitor.results}
        config={{ barWidth: 12, heightClassName: "h-8", showNow: false }}
      />
    </Flex>
  );
};
