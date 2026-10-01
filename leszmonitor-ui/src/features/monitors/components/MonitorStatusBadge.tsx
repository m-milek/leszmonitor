import type { ReactNode } from "react";
import { cn } from "cn";
import { Badge } from "@/components/ui/badge";
import type { Status } from "@/components/common/StatusDot";
import {
  BADGE_SIZE_CLASS,
  type BadgeSize,
} from "@/components/common/badge-size";
import { monitorStatusToStatusDot } from "@/features/monitors/status";
import type { MonitorStatus } from "@/features/monitors/types";

const STATUS_CLASS: Record<Status, string> = {
  up: "bg-lm-status-up text-lm-status-up-foreground",
  down: "bg-lm-status-down text-lm-status-down-foreground",
  pending: "bg-lm-status-pending text-lm-status-pending-foreground",
  paused: "bg-lm-status-paused text-lm-status-paused-foreground",
  unknown: "bg-lm-status-unknown text-lm-status-unknown-foreground",
};

export interface MonitorStatusBadgeProps {
  status: MonitorStatus | undefined;
  size?: BadgeSize;
  children: ReactNode;
}

export const MonitorStatusBadge = ({
  status,
  size = "default",
  children,
}: MonitorStatusBadgeProps) => (
  <Badge
    className={cn(
      STATUS_CLASS[monitorStatusToStatusDot(status)],
      BADGE_SIZE_CLASS[size],
    )}
  >
    {children}
  </Badge>
);
