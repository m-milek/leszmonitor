import { cn } from "@/lib/utils";

interface DividerProps {
  direction: "row" | "column";
  className?: string;
}

export function Divider({ direction, className }: DividerProps) {
  return (
    <div
      className={cn(
        "border-border self-stretch",
        direction === "row" ? "w-full border-t" : "border-l",
        className,
      )}
    />
  );
}
