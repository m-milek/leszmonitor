import type { MonitorResult } from "@/features/monitors/types";
import { monitorStatusToStatusDot } from "@/features/monitors/status";
import { failureReasonToLabel } from "@/features/monitors/failure-reason";
import { type ColumnDef } from "@tanstack/table-core";
import { formatDate } from "@/lib/utils";
import { StatusDot } from "@/components/common/StatusDot";
import { ShortId } from "@/components/common/ShortId";
import { DataTable } from "@/components/common/DataTable";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Minus } from "lucide-react";

export interface MonitorResultsCardProps {
  results: MonitorResult[];
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

export const MonitorResultsCard = ({ results }: MonitorResultsCardProps) => {
  return (
    <Card className="min-w-0">
      <CardHeader>
        <CardTitle>Results</CardTitle>
      </CardHeader>
      <CardContent>
        <ScrollArea className="h-96">
          <DataTable data={results} columns={columns} compact />
        </ScrollArea>
      </CardContent>
    </Card>
  );
};
