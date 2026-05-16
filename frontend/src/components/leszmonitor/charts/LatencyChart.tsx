import type { MonitorResult } from "@/lib/types.ts";
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from "recharts";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart.tsx";
import {
  formatTime,
  generateValuesWithInterval,
} from "@/components/leszmonitor/charts/utils.ts";

export interface LatencyChartProps {
  monitorResults: MonitorResult[];
}

const chartConfig = {
  durationMs: {
    label: "Latency (ms)",
  },
};

export const LatencyChart = ({ monitorResults }: LatencyChartProps) => {
  const mappedData = [...monitorResults].sort(
    (a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime(),
  );

  return (
    <ChartContainer config={chartConfig} className="h-full w-full">
      <AreaChart
        accessibilityLayer
        data={mappedData}
        margin={{ left: -30, top: 5, right: 15 }}
      >
        <CartesianGrid
          vertical={false}
          horizontalValues={generateValuesWithInterval(100, 100)}
        />
        <YAxis domain={[0, 500]} />
        <XAxis dataKey="createdAt" minTickGap={20} tickFormatter={formatTime} />
        <ChartTooltip
          cursor={false}
          content={<ChartTooltipContent indicator="line" />}
        />
        <Area
          dataKey="durationMs"
          type="monotone"
          fillOpacity={0.4}
          fill="var(--chart-1)"
          stroke="var(--chart-2)"
          isAnimationActive={false}
        />
        {/*<ChartLegend content={<ChartLegendContent />} />*/}
      </AreaChart>
    </ChartContainer>
  );
};
