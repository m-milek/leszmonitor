import { cn } from "@/lib/utils";

const COLOR_MAP = {
  success: "bg-green-500",
  failure: "bg-red-700",
  pending: "bg-yellow-500",
} as const;

const SIZE_MAP = {
  sm: "size-2",
  default: "size-3",
  lg: "size-4",
  xl: "size-6",
} as const;

export interface StatusDotProps {
  status: "success" | "failure" | "pending";
  size?: "sm" | "default" | "lg" | "xl";
}

export const StatusDot = ({ status, size = "default" }: StatusDotProps) => {
  const color = COLOR_MAP[status] ?? COLOR_MAP.pending;
  const sizeClass = SIZE_MAP[size] ?? SIZE_MAP.default;

  return <span className={cn("inline-block rounded-full", sizeClass, color)} />;
};
