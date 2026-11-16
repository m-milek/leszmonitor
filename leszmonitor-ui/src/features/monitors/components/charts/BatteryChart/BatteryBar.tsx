import { useState } from "react";
import type { BatteryChartResult } from "@/features/monitors/components/charts/BatteryChart/prepare-results";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Flex } from "@/components/common/Flex";
import { STATUS_BG_CLASS } from "@/components/common/StatusDot";
import { monitorStatusToStatusDot } from "@/features/monitors/status";
import { cn } from "cn";

export const BatteryBar = <T extends BatteryChartResult>({
  result,
  animate,
  width,
  formatResult,
}: {
  result?: T;
  animate?: boolean;
  width: number;
  formatResult: (result: T) => Record<string, string>;
}) => {
  const [animateOnMount] = useState(animate);
  if (!result) {
    return (
      <div className="h-full shrink-0 p-0.5" style={{ width }}>
        <div className={cn("h-full rounded-full", STATUS_BG_CLASS.unknown)} />
      </div>
    );
  }

  const colorClass = STATUS_BG_CLASS[monitorStatusToStatusDot(result.status)];

  return (
    <Popover>
      <PopoverTrigger
        className="h-full shrink-0 cursor-pointer p-0.5 transition-opacity hover:opacity-50"
        style={{ width }}
      >
        <div
          className={cn(
            "h-full rounded-full",
            colorClass,
            animateOnMount &&
              "animate-in zoom-in-0 fade-in slide-in-from-right-2 duration-200 ease-out",
          )}
        />
      </PopoverTrigger>
      <PopoverContent className="w-auto text-sm">
        <Flex direction="column" className="gap-2">
          {Object.entries(formatResult(result)).map(([key, value]) => (
            <Flex key={key} direction="row" className="justify-between gap-4">
              <span className="font-semibold">{key}</span>
              <pre className="font-mono">{value}</pre>
            </Flex>
          ))}
        </Flex>
      </PopoverContent>
    </Popover>
  );
};
