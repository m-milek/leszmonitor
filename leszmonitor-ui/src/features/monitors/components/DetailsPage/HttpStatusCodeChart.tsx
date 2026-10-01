import type { ChartConfig } from "@/components/ui/chart";
import { PieChart } from "@/features/monitors/components/charts/PieChart/PieChartLazy.tsx";

export interface HttpStatusCodeChartProps {
  data: Record<string, number>;
  expectedStatusCodes: number[];
}

const shadedConfig = (statusCodes: string[], colorVar: string): ChartConfig => {
  const step = Math.min(0.16, 0.32 / Math.max(statusCodes.length - 1, 1));
  return Object.fromEntries(
    statusCodes.map((statusCode, index) => [
      statusCode,
      {
        label: statusCode,
        color: `oklch(from var(${colorVar}) calc(l - ${step * index}) c h)`,
      },
    ]),
  );
};

export const HttpStatusCodeChart = ({
  data,
  expectedStatusCodes,
}: HttpStatusCodeChartProps) => {
  const chartData = Object.entries(data)
    .map(([statusCode, count]) => ({ statusCode, count }))
    .sort((a, b) => b.count - a.count);

  const statusCodes = Object.keys(data).sort();
  const isExpected = (statusCode: string) =>
    expectedStatusCodes.includes(Number(statusCode));

  const config: ChartConfig = {
    ...shadedConfig(statusCodes.filter(isExpected), "--lm-status-up"),
    ...shadedConfig(
      statusCodes.filter((statusCode) => !isExpected(statusCode)),
      "--lm-status-down",
    ),
  };

  return (
    <div className="w-full h-64">
      <PieChart
        data={chartData}
        config={config}
        nameKey="statusCode"
        valueKey="count"
      />
    </div>
  );
};
