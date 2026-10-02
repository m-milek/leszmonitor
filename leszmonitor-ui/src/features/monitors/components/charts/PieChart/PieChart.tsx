import { Cell, Pie, PieChart as RechartsPieChart } from "recharts";
import {
  type ChartConfig,
  ChartContainer,
  ChartLegend,
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
          content={() => (
            <div className="flex flex-col items-start gap-1 pl-4">
              {data.map((item) => {
                const name = String(item[nameKey]);
                return (
                  <div key={name} className="flex items-center gap-1.5">
                    <div
                      className="size-2 shrink-0 rounded-[2px]"
                      style={{ backgroundColor: `var(--color-${name})` }}
                    />
                    {config[name]?.label ?? name}
                    <span className="text-muted-foreground tabular-nums">
                      {String(item[valueKey])}
                    </span>
                  </div>
                );
              })}
            </div>
          )}
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
