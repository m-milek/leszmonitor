import type { MonitorStatus } from "@/features/monitors/model/types.ts";
import type { StatusDotProps } from "@/components/common/StatusDot.tsx";

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
