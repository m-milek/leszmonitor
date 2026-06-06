import { type ColumnDef, getCoreRowModel } from "@tanstack/table-core";
import type { AuditLogEntry } from "@/lib/types.ts";
import { flexRender, useReactTable } from "@tanstack/react-table";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx";
import { Badge } from "@/components/ui/badge.tsx";
import { formatDate } from "@/lib/utils.ts";
import { CheckCircle2, LucideDiff, XCircle } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx";
import { ShortId } from "@/components/leszmonitor/tables/audit-log/ShortId.tsx";
import { ResourceDiff } from "@/components/leszmonitor/tables/audit-log/ResourceDiff.tsx";
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
          className="h-6 w-6 text-emerald-500"
          aria-label="Success"
        />
      ) : (
        <XCircle className="h-6 w-6 text-destructive" aria-label="Failed" />
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
    accessorKey: "projectId",
    header: "Project ID",
    cell: ({ row }) => <ShortId value={row.original.projectId} />,
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
          <DialogTrigger asChild>
            <Button variant="ghost">
              <LucideDiff />
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-[90vw] sm:max-w-[1200px] w-full max-h-[90vh] overflow-y-auto">
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
  const table = useReactTable({
    data: entries,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="rounded-md border border-border bg-card">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id} className="hover:bg-transparent">
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id} className="h-12 px-6">
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow
                key={row.id}
                data-state={row.getIsSelected() && "selected"}
                className="transition-colors hover:bg-muted/40"
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id} className="px-6 py-5">
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell
                colSpan={columns.length}
                className="h-32 px-6 text-center"
              >
                No audit log entries yet.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
};
