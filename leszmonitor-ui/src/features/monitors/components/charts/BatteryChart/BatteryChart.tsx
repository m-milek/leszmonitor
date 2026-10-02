import { useEffect, useMemo, useRef, useState } from "react";
import type { MonitorResult } from "@/features/monitors/types";
import { BatteryBar } from "@/features/monitors/components/charts/BatteryChart/BatteryBar";
import { prepareResults } from "@/features/monitors/components/charts/BatteryChart/prepare-results";
import { Flex } from "@/components/common/Flex.tsx";
import { formatDuration } from "@/lib/utils.ts";
import { differenceInSeconds } from "date-fns";

export const BAR_WIDTH = 16;

export const BatteryChart = ({
  monitorResults,
}: {
  monitorResults: MonitorResult[];
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [length, setLength] = useState(0);
  const lastSeenAtRef = useRef<number | null>(null);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;

    const observer = new ResizeObserver(() =>
      setLength(Math.floor(el.clientWidth / BAR_WIDTH)),
    );

    observer.observe(el);
    return () => observer.disconnect();
  }, []);

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
      <div className="flex h-10 w-full overflow-hidden" ref={containerRef}>
        {displayResults.map((res, i) => (
          <BatteryBar
            key={res?.id ?? `empty-${i}`}
            result={res}
            animate={
              res !== undefined &&
              lastSeenAtRef.current !== null &&
              res.createdAt.getTime() > lastSeenAtRef.current
            }
          />
        ))}
      </div>
      <Flex
        direction="row"
        className="justify-between text-muted-foreground"
        style={{ width: displayResults.length * BAR_WIDTH }}
      >
        <span>{formatDuration(oldestResultUntilNowSeconds)} ago</span>
        <span>now</span>
      </Flex>
    </Flex>
  );
};
