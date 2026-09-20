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
import { cn } from "cn";

export interface DataTableProps<T> {
  data: T[];
  columns: ColumnDef<T>[];
  wrapperClassName?: string;
  headRowClassName?: string;
  headClassName?: string;
  bodyRowClassName?: string;
  cellClassName?: string;
  emptyMessage?: string;
  emptyClassName?: string;
}

export const DataTable = <T,>({
  data,
  columns,
  wrapperClassName,
  headRowClassName,
  headClassName,
  bodyRowClassName,
  cellClassName,
  emptyMessage = "No results.",
  emptyClassName = "h-24 text-center",
}: DataTableProps<T>) => {
  const table = useReactTable({
    data: data || [],
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  const table_ = (
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

  if (!wrapperClassName) {
    return table_;
  }

  return <div className={cn(wrapperClassName)}>{table_}</div>;
};
