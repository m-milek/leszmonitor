import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { flexRender, useReactTable } from "@tanstack/react-table";
import { type ColumnDef, getCoreRowModel } from "@tanstack/table-core";

export interface DataTableProps<T> {
  data: T[];
  columns: ColumnDef<T>[];
  emptyMessage?: string;
}

const headRowClassName = "hover:bg-transparent";
const headClassName = "h-12 px-6";
const bodyRowClassName = "transition-colors hover:bg-muted/40";
const cellClassName = "px-6 py-5";
const emptyClassName = "h-32 px-6 text-center";

export const DataTable = <T,>({
  data,
  columns,
  emptyMessage = "No results.",
}: DataTableProps<T>) => {
  const table = useReactTable({
    data: data || [],
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <Table>
      <TableHeader>
        {table.getHeaderGroups().map((headerGroup) => (
          <TableRow key={headerGroup.id} className={headRowClassName}>
            {headerGroup.headers.map((header) => {
              return (
                <TableHead key={header.id} className={headClassName}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                </TableHead>
              );
            })}
          </TableRow>
        ))}
      </TableHeader>
      <TableBody>
        {table.getRowModel().rows?.length ? (
          table.getRowModel().rows.map((row) => (
            <TableRow
              key={row.id}
              data-state={row.getIsSelected() && "selected"}
              className={bodyRowClassName}
            >
              {row.getVisibleCells().map((cell) => (
                <TableCell key={cell.id} className={cellClassName}>
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </TableCell>
              ))}
            </TableRow>
          ))
        ) : (
          <TableRow>
            <TableCell colSpan={columns.length} className={emptyClassName}>
              {emptyMessage}
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
};
