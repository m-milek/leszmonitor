import type { MonitorResult } from "@/features/monitors/types";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Flex } from "@/components/common/Flex";
import { BAR_WIDTH } from "@/features/monitors/components/charts/BatteryChart/BatteryChart";
import { formatResultData } from "@/features/monitors/components/charts/BatteryChart/format-result-data";
import { STATUS_BG_CLASS } from "@/components/common/StatusDot";
import { monitorStatusToStatusDot } from "@/features/monitors/status";
import { cn } from "cn";

export const BatteryBar = ({ result }: { result?: MonitorResult }) => {
  const colorClass =
    STATUS_BG_CLASS[
      result ? monitorStatusToStatusDot(result.status) : "unknown"
    ];

  return (
    <Popover>
      <PopoverTrigger
        className="h-full shrink-0 cursor-pointer transition-opacity hover:opacity-50"
        style={{ width: BAR_WIDTH }}
      >
        <div className={cn("h-full shrink-0 rounded-full m-0.5", colorClass)} />
      </PopoverTrigger>
      <PopoverContent className="w-auto text-sm">
        <Flex direction="column" className="gap-2">
          {result
            ? Object.entries(formatResultData(result)).map(([key, value]) => (
                <Flex
                  key={key}
                  direction="row"
                  className="justify-between gap-4"
                >
                  <span className="font-semibold">{key}</span>
                  <pre className="font-mono">{value}</pre>
                </Flex>
              ))
            : "No data"}
        </Flex>
      </PopoverContent>
    </Popover>
  );
};
