import { lazy, Suspense } from "react";
import type {
  LineChart as LineChartImpl,
  LineChartProps,
} from "./LineChart.tsx";

const LineChartInner = lazy(() =>
  import("./LineChart.tsx").then((m) => ({ default: m.LineChart })),
) as typeof LineChartImpl;

export type { LineChartProps };

export function LineChart<T>(props: LineChartProps<T>) {
  return (
    <Suspense
      fallback={
        <div className="h-full w-full animate-pulse rounded-md bg-muted" />
      }
    >
      <LineChartInner {...props} />
    </Suspense>
  );
}
