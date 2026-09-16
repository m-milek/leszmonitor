import type { MonitorStatus } from "@/features/monitors/types";
import type { StatusDotProps } from "@/components/common/StatusDot";

export const monitorStatusToStatusDot = (
  status: MonitorStatus | undefined,
): StatusDotProps["status"] => {
  switch (status) {
    case "up":
      return "success";
    case "down":
      return "failure";
    default:
      return "pending";
  }
};
