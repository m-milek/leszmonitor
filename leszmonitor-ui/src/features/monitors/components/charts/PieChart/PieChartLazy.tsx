import { lazy, Suspense } from "react";
import type { PieChart as PieChartImpl, PieChartProps } from "./PieChart";

const PieChartInner = lazy(() =>
  import("./PieChart").then((m) => ({ default: m.PieChart })),
) as typeof PieChartImpl;

export type { PieChartProps };

export function PieChart<T>(props: PieChartProps<T>) {
  return (
    <Suspense fallback={null}>
      <PieChartInner {...props} />
    </Suspense>
  );
}
