import { Cell, Pie, PieChart as RechartsPieChart } from "recharts";
import {
  type ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { CHART_CONFIG } from "@/features/monitors/components/charts/charts-config";

export interface PieChartProps<T> {
  data: T[];
  config: ChartConfig;
  nameKey: Extract<keyof T, string>;
  valueKey: Extract<keyof T, string>;
}

export function PieChart<T>({
  data,
  config,
  nameKey,
  valueKey,
}: PieChartProps<T>) {
  return (
    <ChartContainer config={config} className="h-full w-full">
      <RechartsPieChart accessibilityLayer>
        <ChartTooltip
          cursor={false}
          content={<ChartTooltipContent hideLabel />}
        />
        <ChartLegend
          layout="vertical"
          align="right"
          verticalAlign="middle"
          content={
            <ChartLegendContent
              nameKey={nameKey}
              className="flex-col items-start gap-1 pt-0 pl-4"
            />
          }
        />
        <Pie
          data={data}
          dataKey={valueKey}
          nameKey={nameKey}
          animationDuration={CHART_CONFIG.ANIMATION.duration}
          animationEasing={CHART_CONFIG.ANIMATION.easing}
        >
          {data.map((item) => (
            <Cell
              key={String(item[nameKey])}
              fill={`var(--color-${String(item[nameKey])})`}
            />
          ))}
        </Pie>
      </RechartsPieChart>
    </ChartContainer>
  );
}
