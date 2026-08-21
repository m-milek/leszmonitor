import type { MonitorResult } from "@/lib/types.ts";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { BAR_WIDTH } from "@/components/leszmonitor/charts/BatteryChart/BatteryChart.tsx";
import { formatResultData } from "@/components/leszmonitor/charts/BatteryChart/format-result-data.ts";
import { cn } from "@/lib/utils.ts";

const STATUS_COLORS: Record<string, string> = {
  up: "bg-green-500",
  down: "bg-red-700",
  paused: "bg-blue-500",
  maintenance: "bg-yellow-500",
  default: "bg-gray-300",
};

export const BatteryBar = ({ result }: { result?: MonitorResult }) => {
  const colorClass = result
    ? STATUS_COLORS[result.status] || STATUS_COLORS.default
    : STATUS_COLORS.default;

  return (
    <Popover>
      <PopoverTrigger
        className={`h-full shrink-0 cursor-pointer hover:opacity-50 transition-opacity`}
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
