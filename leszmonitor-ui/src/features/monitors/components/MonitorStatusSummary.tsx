import { cn } from "cn";
import { Button } from "@/components/ui/button";
import { StatusDot } from "@/components/common/StatusDot";
import {
  getEffectiveStatus,
  monitorStatusToStatusDot,
} from "@/features/monitors/status";
import {
  monitorStatuses,
  type MonitorWithStatus,
} from "@/features/monitors/types";

export interface MonitorStatusSummaryProps {
  monitors: MonitorWithStatus[];
  value: string[];
  onChange: (value: string[]) => void;
  className?: string;
}

export const MonitorStatusSummary = ({
  monitors,
  value,
  onChange,
  className,
}: MonitorStatusSummaryProps) => {
  const counts = monitorStatuses
    .map((status) => ({
      status,
      count: monitors.filter((m) => getEffectiveStatus(m) === status).length,
    }))
    .filter(({ count }) => count > 0);

  const toggle = (status: string) =>
    onChange(
      value.includes(status)
        ? value.filter((s) => s !== status)
        : [...value, status],
    );

  return (
    <div className={cn("flex flex-wrap items-center gap-1", className)}>
      {counts.map(({ status, count }) => (
        <Button
          key={status}
          variant={value.includes(status) ? "secondary" : "ghost"}
          aria-pressed={value.includes(status)}
          onClick={() => toggle(status)}
        >
          <StatusDot status={monitorStatusToStatusDot(status)} size="sm" />
          <span className="font-semibold tabular-nums">{count}</span>
          <span className="text-muted-foreground">{status}</span>
        </Button>
      ))}
    </div>
  );
};
