import type { ChartConfig } from "@/components/ui/chart";
import { PieChart } from "@/features/monitors/components/charts/PieChart/PieChartLazy.tsx";

export interface HttpStatusCodeChartProps {
  data: Record<string, number>;
}

const getStatusCodeColor = (statusCode: string) => {
  switch (statusCode[0]) {
    case "2":
      return "var(--lm-status-up)";
    case "3":
      return "var(--lm-status-paused)";
    case "4":
      return "var(--lm-status-pending)";
    case "5":
      return "var(--lm-status-down)";
    default:
      return "var(--lm-status-unknown)";
  }
};

export const HttpStatusCodeChart = ({ data }: HttpStatusCodeChartProps) => {
  const chartData = Object.entries(data)
    .map(([statusCode, count]) => ({ statusCode, count }))
    .sort((a, b) => b.count - a.count);

  const config: ChartConfig = Object.fromEntries(
    chartData.map(({ statusCode }) => [
      statusCode,
      { label: statusCode, color: getStatusCodeColor(statusCode) },
    ]),
  );

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
