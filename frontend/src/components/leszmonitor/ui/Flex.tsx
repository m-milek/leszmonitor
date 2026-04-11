import type { CSSProperties, ReactNode } from "react";
import { cn } from "@/lib/utils.ts";

export interface FlexProps {
  children: ReactNode;
  direction?: "horizontal" | "vertical";
  align?: CSSProperties["alignItems"];
  gap?: CSSProperties["gap"];
  justify?: CSSProperties["justifyContent"];
  wrap?: CSSProperties["flexWrap"];
  className?: string;
}

export function Flex({
  children,
  direction = "horizontal",
  align = "center",
  gap = "1rem",
  justify = "flex-start",
  wrap = "wrap",
  className,
}: FlexProps) {
  return (
    <div
      className={cn(
        "flex",
        direction === "vertical" ? "flex-col" : "flex-row",
        className,
      )}
      style={{
        alignItems: align,
        justifyContent: justify,
        gap,
        flexWrap: wrap,
      }}
    >
      {children}
    </div>
  );
}
