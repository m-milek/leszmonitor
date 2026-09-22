import type { MonitorStatus } from "@/features/monitors/types";
import type { Status } from "@/components/common/StatusDot";

export const monitorStatusToStatusDot = (
  status: MonitorStatus | undefined,
): Status => {
  switch (status) {
    case "up":
      return "up";
    case "down":
      return "down";
    case "paused":
      return "paused";
    default:
      return "pending";
  }
};
