import { useEffect, useMemo, useRef, useState } from "react";
import type { MonitorResult } from "@/features/monitors/types";
import { Flex } from "@/components/common/Flex";
import { BatteryBar } from "@/features/monitors/components/charts/BatteryChart/BatteryBar";
import { prepareResults } from "@/features/monitors/components/charts/BatteryChart/prepare-results";

export interface BatteryChartProps {
  length?: number;
  monitorResults: MonitorResult[];
}

export const BAR_WIDTH = 16;
const GAP = 0;
const TOTAL_WIDTH = BAR_WIDTH + GAP;

export const BatteryChart = ({
  length: defaultLength,
  monitorResults,
}: BatteryChartProps) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [length, setLength] = useState(defaultLength ?? 50);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;

    const observer = new ResizeObserver(([entry]) => {
      const width = entry?.contentRect.width ?? el.clientWidth;
      if (width > 0) {
        setLength(Math.max(1, Math.floor((width + GAP) / TOTAL_WIDTH)));
      }
    });

    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  const displayResults = useMemo(() => {
    return prepareResults(monitorResults, length);
  }, [monitorResults, length]);

  return (
    <div className="w-full h-8 relative" ref={containerRef}>
      <Flex
        direction="row"
        className="justify-start items-center absolute inset-0"
        style={{ gap: GAP }}
      >
        {displayResults.map((res, i) => (
          <BatteryBar key={i} result={res} />
        ))}
      </Flex>
    </div>
  );
};
