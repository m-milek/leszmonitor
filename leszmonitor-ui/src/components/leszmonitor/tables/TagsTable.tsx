import { type ColumnDef, getCoreRowModel } from "@tanstack/table-core";
import { Minus } from "lucide-react";
import { flexRender, useReactTable } from "@tanstack/react-table";
import type { Tag as TagModel } from "@/lib/types.ts";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx";
import { Tag } from "@/components/leszmonitor/Tag.tsx";
import { ShortId } from "@/components/leszmonitor/tables/audit-log/ShortId.tsx";
import { formatDate } from "@/lib/utils.ts";
import { normalizeHexColor } from "@/lib/tagColors.ts";
import { DeleteTagDialog } from "@/components/leszmonitor/dialogs/DeleteTagDialog.tsx";

export interface TagsTableProps {
  tags: TagModel[];
}

const columns: ColumnDef<TagModel>[] = [
  {
    accessorKey: "name",
    header: "Tag",
    cell: ({ row }) => <Tag tag={row.original} />,
  },
  {
    accessorKey: "description",
    header: "Description",
    cell: ({ row }) =>
      row.original.description ? (
        <span>{row.original.description}</span>
      ) : (
        <Minus
          className="h-4 w-4 text-muted-foreground"
          aria-label="No description"
        />
      ),
  },
  {
    accessorKey: "colorHex",
    header: "Color",
    cell: ({ row }) => (
      <code className="text-muted-foreground uppercase">
        {normalizeHexColor(row.original.colorHex)}
      </code>
    ),
  },
  {
    accessorKey: "id",
    header: "ID",
    cell: ({ row }) => <ShortId value={row.original.id} />,
  },
  {
    accessorKey: "createdAt",
    header: "Created",
    cell: ({ row }) => (
      <span className="whitespace-nowrap">
        {formatDate(row.original.createdAt)}
      </span>
    ),
  },
  {
    id: "actions",
    header: "",
    cell: ({ row }) => (
      <div className="flex justify-end">
        <DeleteTagDialog tag={row.original} />
      </div>
    ),
  },
];

export const TagsTable = ({ tags }: TagsTableProps) => {
  const table = useReactTable({
    data: tags,
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
                No tags yet.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
};
