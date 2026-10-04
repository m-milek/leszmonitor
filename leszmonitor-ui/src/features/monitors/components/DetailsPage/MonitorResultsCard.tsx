import { type MonitorResult, monitorStatuses } from "@/features/monitors/types";
import { monitorStatusToStatusDot } from "@/features/monitors/status";
import { failureReasonToLabel } from "@/features/monitors/failure-reason";
import { type ColumnDef } from "@tanstack/table-core";
import { formatDate } from "@/lib/utils";
import { StatusDot } from "@/components/common/StatusDot";
import { ShortId } from "@/components/common/ShortId";
import { DataTable } from "@/components/common/DataTable";
import {
  Card,
  CardAction,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { EyeIcon, LucideMinus, Minus } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useMemo, useState } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { LMFacetedFilter } from "@/components/form/LMFacetedFilter";

export interface MonitorResultsCardProps {
  results?: MonitorResult[];
}

const formatDetails = (result: MonitorResult) => {
  switch (result.monitorType) {
    case "http":
      return `HTTP ${result.details.statusCode}`;
    case "tcp":
      return `${result.details.latencyMs} ms latency`;
    case "dns":
      return `${result.details.resolvedRecords?.length ?? 0} records`;
    case "push":
      return `Message: ${result.details.rawMessage?.length ?? 0} bytes`;
  }
};

interface RawResultDialogProps {
  result: MonitorResult | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

const RawResultDialog = ({
  result,
  open,
  onOpenChange,
}: RawResultDialogProps) => (
  <Dialog open={open} onOpenChange={onOpenChange}>
    <DialogContent className="sm:max-w-2xl">
      <DialogHeader>
        <DialogTitle>Result {result?.id}</DialogTitle>
      </DialogHeader>
      <pre className="max-h-[70vh] overflow-auto font-mono text-xs">
        {JSON.stringify(result, null, 2)}
      </pre>
    </DialogContent>
  </Dialog>
);

const baseColumns: ColumnDef<MonitorResult>[] = [
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
    cell: ({ row }) =>
      row.original.durationMs ? (
        `${row.original.durationMs} ms`
      ) : (
        <LucideMinus />
      ),
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
  const [selectedResult, setSelectedResult] = useState<MonitorResult | null>(
    null,
  );
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [statusFilter, setStatusFilter] = useState<string[]>([]);

  const filteredResults = statusFilter.length
    ? results?.filter((result) => statusFilter.includes(result.status))
    : results;

  const columns = useMemo<ColumnDef<MonitorResult>[]>(
    () => [
      ...baseColumns,
      {
        id: "raw",
        header: "",
        cell: ({ row }) => (
          <Button
            variant="ghost"
            size="icon-sm"
            aria-label="Show raw JSON"
            onClick={() => {
              setSelectedResult(row.original);
              setIsDialogOpen(true);
            }}
          >
            <EyeIcon />
          </Button>
        ),
      },
    ],
    [],
  );

  return (
    <Card className="min-w-0">
      <CardHeader>
        <CardTitle>Results</CardTitle>
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
        {filteredResults ? (
          <DataTable
            data={filteredResults}
            columns={columns}
            compact
            pageSize={10}
            getRowId={(result) => result.id}
          />
        ) : (
          <div className="flex flex-col gap-2">
            {Array.from({ length: 10 }, (_, i) => (
              <Skeleton key={i} className="h-8 w-full" />
            ))}
          </div>
        )}
        <RawResultDialog
          result={selectedResult}
          open={isDialogOpen}
          onOpenChange={setIsDialogOpen}
        />
      </CardContent>
    </Card>
  );
};
