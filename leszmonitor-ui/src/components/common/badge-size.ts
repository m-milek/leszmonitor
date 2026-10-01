export type BadgeSize = "default" | "lg";

export const BADGE_SIZE_CLASS: Record<BadgeSize, string> = {
  default: "",
  lg: "h-(--text-3xl) text-sm",
};
