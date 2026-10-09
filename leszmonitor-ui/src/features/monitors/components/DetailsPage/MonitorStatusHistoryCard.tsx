import {
  type MonitorStatusPeriod,
  monitorStatuses,
} from "@/features/monitors/types";
import { type ColumnDef } from "@tanstack/react-table";
import { formatDate, formatDuration } from "@/lib/utils";
import { StatusDot } from "@/components/common/StatusDot";
import { monitorStatusToStatusDot } from "@/features/monitors/status";
import { DataTable, dataTableFeatures } from "@/components/common/DataTable";
import {
  Card,
  CardAction,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useMemo, useState } from "react";
import { LMFacetedFilter } from "@/components/form/LMFacetedFilter";
import { Skeleton } from "@/components/ui/skeleton";

export interface MonitorStatusHistoryCardProps {
  statusHistory?: MonitorStatusPeriod[];
}

const baseColumns: ColumnDef<typeof dataTableFeatures, MonitorStatusPeriod>[] =
  [
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
      accessorKey: "startedAt",
      header: "Started",
      cell: ({ row }) => (
        <span className="whitespace-nowrap">
          {formatDate(row.original.startedAt)}
        </span>
      ),
    },
    {
      accessorKey: "endedAt",
      header: "Ended",
      cell: ({ row }) => (
        <span className="whitespace-nowrap">
          {row.original.endedAt ? formatDate(row.original.endedAt) : "Current"}
        </span>
      ),
    },
    {
      accessorKey: "durationSeconds",
      header: "Duration",
      cell: ({ row }) => (
        <span className="whitespace-nowrap">
          {formatDuration(row.original.durationSeconds)}
        </span>
      ),
    },
  ];

export const MonitorStatusHistoryCard = ({
  statusHistory,
}: MonitorStatusHistoryCardProps) => {
  const columns = useMemo<
    ColumnDef<typeof dataTableFeatures, MonitorStatusPeriod>[]
  >(() => baseColumns, []);
  const [statusFilter, setStatusFilter] = useState<string[]>([]);

  const filteredStatusHistory = statusFilter.length
    ? statusHistory?.filter((period) => statusFilter.includes(period.status))
    : statusHistory;

  return (
    <Card className="min-w-0">
      <CardHeader>
        <CardTitle>Status History</CardTitle>
        <CardAction>
          <LMFacetedFilter
            title="Status"
            options={monitorStatuses.map((status) => ({
              value: status,
              label: status.toUpperCase(),
            }))}
            value={statusFilter}
            onChange={setStatusFilter}
          />
        </CardAction>
      </CardHeader>
      <CardContent>
        {filteredStatusHistory ? (
          <DataTable
            data={filteredStatusHistory}
            columns={columns}
            compact
            pageSize={10}
          />
        ) : (
          <div className="flex flex-col gap-2">
            {Array.from({ length: 10 }, (_, i) => (
              <Skeleton key={i} className="h-8 w-full" />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
};
