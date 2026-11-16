import { useEffect, useMemo, useRef, useState } from "react";
import { BatteryBar } from "@/features/monitors/components/charts/BatteryChart/BatteryBar";
import {
  type BatteryChartResult,
  prepareResults,
} from "@/features/monitors/components/charts/BatteryChart/prepare-results";
import { cn } from "cn";
import { Flex } from "@/components/common/Flex";
import { formatDuration } from "@/lib/utils";
import { differenceInSeconds } from "date-fns";

export interface BatteryChartConfig {
  barWidth: number;
  heightClassName: string;
  showTimeRange: boolean;
  showNow: boolean;
}

const DEFAULT_CONFIG: BatteryChartConfig = {
  barWidth: 16,
  heightClassName: "h-10",
  showTimeRange: true,
  showNow: true,
};

export interface BatteryChartProps<T extends BatteryChartResult> {
  monitorResults: T[];
  config?: Partial<BatteryChartConfig>;
  formatResult?: (result: T) => Record<string, string>;
}

const defaultFormatResult = (
  result: BatteryChartResult,
): Record<string, string> => ({
  Status: result.status,
  "Created At": result.createdAt.toLocaleString(),
});

export const BatteryChart = <T extends BatteryChartResult>({
  monitorResults,
  config,
  formatResult = defaultFormatResult,
}: BatteryChartProps<T>) => {
  const { barWidth, heightClassName, showTimeRange, showNow } = {
    ...DEFAULT_CONFIG,
    ...config,
  };
  const containerRef = useRef<HTMLDivElement>(null);
  const [length, setLength] = useState(0);
  const lastSeenAtRef = useRef<number | null>(null);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;

    const observer = new ResizeObserver(() =>
      setLength(Math.floor(el.clientWidth / barWidth)),
    );

    observer.observe(el);
    return () => observer.disconnect();
  }, [barWidth]);

  const displayResults = useMemo(() => {
    return prepareResults(monitorResults, length);
  }, [monitorResults, length]);

  const newestAt = displayResults.at(-1)?.createdAt.getTime() ?? null;

  useEffect(() => {
    lastSeenAtRef.current = newestAt;
  }, [newestAt]);

  const oldestAt = displayResults.at(0)?.createdAt;

  return (
    <Flex direction="column">
      <div
        className={cn("flex w-full overflow-hidden", heightClassName)}
        ref={containerRef}
      >
        {displayResults.map((res, i) => (
          <BatteryBar
            key={res?.createdAt.toISOString() ?? `empty-${i}`}
            result={res}
            formatResult={formatResult}
            width={barWidth}
            animate={
              res !== undefined &&
              lastSeenAtRef.current !== null &&
              res.createdAt.getTime() > lastSeenAtRef.current
            }
          />
        ))}
      </div>
      {showTimeRange && (
        <Flex
          direction="row"
          className="justify-between text-muted-foreground"
          style={{ width: displayResults.length * barWidth }}
        >
          {oldestAt && (
            <span>
              {formatDuration(differenceInSeconds(new Date(), oldestAt))} ago
            </span>
          )}
          {showNow && <span className="ml-auto">now</span>}
        </Flex>
      )}
    </Flex>
  );
};
