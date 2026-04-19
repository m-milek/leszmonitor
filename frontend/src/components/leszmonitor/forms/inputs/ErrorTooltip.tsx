import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip.tsx";
import type { ReactNode } from "react";

interface ErrorTooltipProps {
  isOpen: boolean;
  message: string;
  children: ReactNode;
  side?: "top" | "right" | "bottom" | "left";
}

export function ErrorTooltip({
  isOpen,
  message,
  children,
  side = "top",
}: ErrorTooltipProps) {
  return (
    <Tooltip open={isOpen}>
      <TooltipTrigger asChild>{children}</TooltipTrigger>
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
