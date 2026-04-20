import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip.tsx";
import type { ReactNode } from "react";

interface ErrorTooltipProps {
  isOpen?: boolean;
  message?: string;
  side?: "top" | "right" | "bottom" | "left";
  children: ReactNode;
}

export function ErrorTooltip({
  isOpen = false,
  message = "",
  side = "top",
  children,
}: ErrorTooltipProps) {
  return (
    <Tooltip open={isOpen}>
      <TooltipTrigger asChild>
        <span>{children}</span>
      </TooltipTrigger>
      {isOpen && (
        <TooltipContent
          side={side}
          className="bg-destructive text-white border-destructive"
          arrowClassName="bg-destructive fill-destructive"
        >
          {message}
        </TooltipContent>
      )}
    </Tooltip>
  );
}
