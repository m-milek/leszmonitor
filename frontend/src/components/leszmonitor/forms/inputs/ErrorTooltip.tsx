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
  side = "bottom",
}: ErrorTooltipProps) {
  return (
    <Tooltip open={isOpen}>
      <TooltipTrigger asChild>{children}</TooltipTrigger>
      {isOpen && (
        <TooltipContent
          side={side}
          className="bg-red-500 text-white border-red-600"
          arrowClassName="bg-red-500 fill-red-500"
        >
          {message}
        </TooltipContent>
      )}
    </Tooltip>
  );
}
