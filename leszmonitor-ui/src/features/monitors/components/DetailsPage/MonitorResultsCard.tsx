import { MonitorsApi } from "@/features/monitors/monitors-api";
import { monitorStatusToStatusDot } from "@/features/monitors/status";
import { failureReasonToLabel } from "@/features/monitors/failure-reason";
import type { Monitor, MonitorResult } from "@/features/monitors/types";
import type { Pagination } from "@/lib/types";
import { useQuery } from "@tanstack/react-query";
import { type ColumnDef } from "@tanstack/table-core";
import { QUERY_KEYS } from "@/lib/consts";
import { formatDate } from "@/lib/utils";
import { StatusDot } from "@/components/common/StatusDot";
import { ShortId } from "@/components/common/ShortId";
import { DataTable } from "@/components/common/DataTable";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Minus } from "lucide-react";

export interface MonitorResultsCardProps {
  monitor: Monitor;
  pagination: Pagination;
}

const formatDetails = (result: MonitorResult) => {
  switch (result.monitorType) {
    case "http":
      return `HTTP ${result.details.statusCode}`;
    case "tcp":
      return `${result.details.latencyMs} ms latency`;
    case "dns":
      return `${result.details.resolvedRecords?.length ?? 0} records`;
  }
};

const columns: ColumnDef<MonitorResult>[] = [
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => (
      <span className="flex items-center gap-2 font-mono">
        <StatusDot status={monitorStatusToStatusDot(row.original.status)} />
        {row.original.status.toUpperCase()}
      </span>
    ),
  },
  {
    accessorKey: "createdAt",
    header: "Time",
    cell: ({ row }) => (
      <span className="whitespace-nowrap">
        {formatDate(row.original.createdAt)}
      </span>
    ),
  },
  {
    accessorKey: "durationMs",
    header: "Duration",
    cell: ({ row }) => `${row.original.durationMs} ms`,
  },
  {
    id: "details",
    header: "Details",
    cell: ({ row }) => formatDetails(row.original),
  },
  {
    accessorKey: "failures",
    header: "Failures",
    cell: ({ row }) =>
      row.original.failures?.length ? (
        row.original.failures
          .map((failure) => failureReasonToLabel(failure.reason))
          .join(", ")
      ) : (
        <Minus
          className="size-4 text-muted-foreground"
          aria-label="No failures"
        />
      ),
  },
  {
    accessorKey: "id",
    header: "ID",
    cell: ({ row }) => <ShortId value={row.original.id} />,
  },
];

export const MonitorResultsCard = ({
  monitor,
  pagination,
}: MonitorResultsCardProps) => {
  const { data: results } = useQuery({
    queryKey: [QUERY_KEYS.MONITOR_RESULTS, monitor.id, pagination],
    queryFn: () => MonitorsApi.results.getPage(monitor.id, pagination),
  });

  return (
    <Card className="min-w-0">
      <CardHeader>
        <CardTitle>Results</CardTitle>
      </CardHeader>
      <CardContent>
        <ScrollArea className="h-96">
          <DataTable data={results ?? []} columns={columns} compact />
        </ScrollArea>
      </CardContent>
    </Card>
  );
};
