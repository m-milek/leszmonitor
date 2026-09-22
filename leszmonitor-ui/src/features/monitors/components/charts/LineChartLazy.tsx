import { lazy, Suspense } from "react";
import type { LineChart as LineChartImpl, LineChartProps } from "./LineChart";

const LineChartInner = lazy(() =>
  import("./LineChart").then((m) => ({ default: m.LineChart })),
) as typeof LineChartImpl;

export type { LineChartProps };

export function LineChart<T>(props: LineChartProps<T>) {
  return (
    <Suspense fallback={null}>
      <LineChartInner {...props} />
    </Suspense>
  );
}
