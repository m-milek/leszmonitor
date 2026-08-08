import {
  Area,
  AreaChart,
  CartesianGrid,
  matchByDataKey,
  XAxis,
  YAxis,
} from "recharts";
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart.tsx";
import { CHART_CONFIG } from "@/components/leszmonitor/charts/charts-config.ts";

export interface LineChartProps<T> {
  data: T[];
  config: ChartConfig;
  timestampExtractor: (item: T) => number;
  xAxisKey: Extract<keyof T, string>;
  yAxisKey: Extract<keyof T, string>;
  uniqueMatchKey: Extract<keyof T, string>;
  xAxisTickFormatter?: (value: any) => string;
  yAxisDomain?: [
    number | "auto" | "dataMin" | "dataMax",
    number | "auto" | "dataMin" | "dataMax",
  ];
}

export function LineChart<T>({
  data,
  config,
  timestampExtractor,
  xAxisKey,
  yAxisKey,
  uniqueMatchKey,
  xAxisTickFormatter,
  yAxisDomain = [0, "auto"],
}: LineChartProps<T>) {
  // Sort data chronologically and trim to the visible window size
  const windowData = [...data]
    .sort((a, b) => timestampExtractor(a) - timestampExtractor(b))
    .slice(-CHART_CONFIG.VISIBLE_POINTS);

  return (
    <ChartContainer config={config} className="h-full w-full">
      <AreaChart
        accessibilityLayer
        data={windowData}
        margin={CHART_CONFIG.MARGIN}
      >
        <CartesianGrid vertical={false} />
        <YAxis domain={yAxisDomain} />
        <XAxis
          dataKey={xAxisKey}
          minTickGap={CHART_CONFIG.X_AXIS.minTickGap}
          tickFormatter={xAxisTickFormatter}
          allowDataOverflow
        />
        <ChartTooltip
          cursor={false}
          content={<ChartTooltipContent indicator="line" />}
        />
        <Area
          dataKey={yAxisKey}
          type="monotone"
          fillOpacity={0.4}
          fill="var(--chart-1)"
          stroke="var(--chart-2)"
          animationDuration={CHART_CONFIG.ANIMATION.duration}
          animationEasing={CHART_CONFIG.ANIMATION.easing}
          animationMatchBy={matchByDataKey(uniqueMatchKey)}
        />
      </AreaChart>
    </ChartContainer>
  );
}
