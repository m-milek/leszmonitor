import { Flex } from "@/components/common/Flex";
import { StatusDot } from "@/components/common/StatusDot";
import { BatteryChart } from "@/features/monitors/components/charts/BatteryChart/BatteryChart";
import { monitorStatusToStatusDot } from "@/features/monitors/status";
import type { PublicMonitor } from "@/features/status-pages/types";
import { formatDuration } from "@/lib/utils";

export interface PublicMonitorRowProps {
  monitor: PublicMonitor;
}

export const PublicMonitorRow = ({ monitor }: PublicMonitorRowProps) => {
  const { latency, statusChange, uptime } = monitor.stats;
  const hasResults = uptime.statusToCount !== null;

  return (
    <Flex direction="column" className="gap-2 py-3">
      <Flex direction="row" className="items-center justify-between gap-4">
        <Flex direction="row" className="items-center gap-2">
          <StatusDot status={monitorStatusToStatusDot(monitor.status)} />
          <span className="font-medium">{monitor.name}</span>
        </Flex>
        <Flex
          direction="row"
          className="text-muted-foreground gap-4 text-sm whitespace-nowrap"
        >
          <span>
            {monitor.status} for{" "}
            {formatDuration(statusChange.secondsInCurrentStatus)}
          </span>
          {hasResults && <span>{Math.round(latency.avg)} ms avg</span>}
          <span>
            {hasResults
              ? `${(uptime.statusToPercentage?.up ?? 0).toFixed(2)}% uptime`
              : "No data"}
          </span>
        </Flex>
      </Flex>
      <BatteryChart
        monitorResults={monitor.results}
        config={{ barWidth: 12, heightClassName: "h-8", showNow: false }}
      />
    </Flex>
  );
};
