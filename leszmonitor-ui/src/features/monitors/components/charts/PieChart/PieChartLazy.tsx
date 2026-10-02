import { lazy, Suspense } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import type { PieChart as PieChartImpl, PieChartProps } from "./PieChart";

const PieChartInner = lazy(() =>
  import("./PieChart").then((m) => ({ default: m.PieChart })),
) as typeof PieChartImpl;

export type { PieChartProps };

export function PieChart<T>(props: PieChartProps<T>) {
  return (
    <Suspense fallback={<Skeleton className="aspect-video h-full" />}>
      <PieChartInner {...props} />
    </Suspense>
  );
}
