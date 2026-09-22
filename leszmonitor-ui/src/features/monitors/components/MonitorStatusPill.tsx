import { Badge } from "@/components/ui/badge";
import type { Monitor, MonitorRunState } from "@/features/monitors/types";

export interface MonitorStatusPillProps {
  monitor: Monitor;
}

const mapMonitorState = (state: MonitorRunState) => {
  switch (state) {
    case "active":
      return {
        text: "Active",
        className: "bg-lm-status-up text-lm-status-up-foreground",
      };
    case "paused":
      return {
        text: "Paused",
        className: "bg-lm-status-paused text-lm-status-paused-foreground",
      };
    default:
      return {
        text: "Invalid",
        className: "bg-lm-status-unknown text-lm-status-unknown-foreground",
      };
  }
};

export const MonitorStatusPill = ({ monitor }: MonitorStatusPillProps) => {
  const { text, className } = mapMonitorState(monitor.runState);

  return <Badge className={className}>{text}</Badge>;
};
