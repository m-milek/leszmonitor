import { cn } from "cn";

export type Status = "up" | "down" | "pending" | "paused" | "unknown";

export const STATUS_BG_CLASS: Record<Status, string> = {
  up: "bg-lm-status-up",
  down: "bg-lm-status-down",
  pending: "bg-lm-status-pending",
  paused: "bg-lm-status-paused",
  unknown: "bg-lm-status-unknown",
};

const SIZE_MAP = {
  sm: "size-2",
  default: "size-3",
  lg: "size-4",
  xl: "size-6",
} as const;

export interface StatusDotProps {
  status: Status;
  size?: "sm" | "default" | "lg" | "xl";
  className?: string;
}

export const StatusDot = ({
  status,
  size = "default",
  className,
}: StatusDotProps) => (
  <span
    className={cn(
      "inline-block rounded-full",
      SIZE_MAP[size],
      STATUS_BG_CLASS[status],
      className,
    )}
  />
);
