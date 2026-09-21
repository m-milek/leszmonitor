import { type ColumnDef } from "@tanstack/table-core";
import type { AuditLogEntry } from "@/features/audit-log/types";
import { DataTable } from "@/components/common/DataTable";
import { Badge } from "@/components/ui/badge";
import { formatDate } from "@/lib/utils";
import { CheckCircle2, LucideDiff, XCircle } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { ShortId } from "@/components/common/ShortId";
import { ResourceDiff } from "@/features/audit-log/components/ResourceDiff";
import { Button } from "@/components/ui/button";

export interface AuditLogTableProps {
  entries: AuditLogEntry[];
}

const columns: ColumnDef<AuditLogEntry>[] = [
  {
    accessorKey: "createdAt",
    header: "Timestamp",
    cell: ({ row }) => (
      <span className="whitespace-nowrap">
        {formatDate(row.original.createdAt)}
      </span>
    ),
  },
  {
    accessorKey: "isSuccess",
    header: "Status",
    cell: ({ row }) =>
      row.original.isSuccess ? (
        <CheckCircle2
          className="size-6 text-lm-status-up"
          aria-label="Success"
        />
      ) : (
        <XCircle className="size-6 text-destructive" aria-label="Failed" />
      ),
  },
  {
    accessorKey: "username",
    header: "User",
  },
  {
    accessorKey: "action",
    header: "Action",
    cell: ({ row }) => <Badge variant="secondary">{row.original.action}</Badge>,
  },
  {
    accessorKey: "resourceId",
    header: "Resource ID",
    cell: ({ row }) => <ShortId value={row.original.resourceId} />,
  },
  {
    header: "Diff",
    cell: ({ row }) =>
      !row.original.before && !row.original.after ? (
        <Button variant="ghost" disabled>
          —
        </Button>
      ) : (
        <Dialog>
          <DialogTrigger render={<Button variant="ghost" />}>
            <LucideDiff />
          </DialogTrigger>
          <DialogContent className="max-w-[90vw] sm:max-w-300 w-full max-h-[90vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>Resource Diff</DialogTitle>
            </DialogHeader>
            <div className="mt-4">
              <ResourceDiff
                before={row.original.before}
                after={row.original.after}
              />
            </div>
          </DialogContent>
        </Dialog>
      ),
  },
  {
    accessorKey: "traceId",
    header: "Trace ID",
    cell: ({ row }) => <ShortId value={row.original.traceId} />,
  },
];

export const AuditLogTable = ({ entries }: AuditLogTableProps) => {
  return (
    <DataTable
      data={entries}
      columns={columns}
      emptyMessage="No audit log entries yet."
    />
  );
};
