import { useEffect, useMemo, useRef, useState } from "react";
import type { MonitorResult } from "@/features/monitors/types";
import { BatteryBar } from "@/features/monitors/components/charts/BatteryChart/BatteryBar";
import { prepareResults } from "@/features/monitors/components/charts/BatteryChart/prepare-results";
import { cn } from "cn";
import { Flex } from "@/components/common/Flex";
import { formatDuration } from "@/lib/utils";
import { differenceInSeconds } from "date-fns";

export interface BatteryChartConfig {
  barWidth: number;
  heightClassName: string;
  showTimeRange: boolean;
}

const DEFAULT_CONFIG: BatteryChartConfig = {
  barWidth: 16,
  heightClassName: "h-10",
  showTimeRange: true,
};

export interface BatteryChartProps {
  monitorResults: MonitorResult[];
  config?: Partial<BatteryChartConfig>;
}

export const BatteryChart = ({ monitorResults, config }: BatteryChartProps) => {
  const { barWidth, heightClassName, showTimeRange } = {
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

  const oldestResultUntilNowSeconds = differenceInSeconds(
    new Date(),
    displayResults.at(0)?.createdAt ?? new Date(),
  );

  return (
    <Flex direction="column">
      <div
        className={cn("flex w-full overflow-hidden", heightClassName)}
        ref={containerRef}
      >
        {displayResults.map((res, i) => (
          <BatteryBar
            key={res?.id ?? `empty-${i}`}
            result={res}
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
          <span>{formatDuration(oldestResultUntilNowSeconds)} ago</span>
          <span>now</span>
        </Flex>
      )}
    </Flex>
  );
};
