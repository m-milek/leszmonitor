import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { flexRender, useReactTable } from "@tanstack/react-table";
import {
  type ColumnDef,
  getCoreRowModel,
  getPaginationRowModel,
} from "@tanstack/table-core";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

export interface DataTableProps<T> {
  data: T[];
  columns: ColumnDef<T>[];
  emptyMessage?: string;
  compact?: boolean;
  pageSize?: number;
  getRowId?: (row: T) => string;
}

const headRowClassName = "hover:bg-transparent";
const headClassName = "h-12 px-6";
const compactHeadClassName = "h-9 px-3";
const bodyRowClassName = "transition-colors hover:bg-muted/40";
const cellClassName = "px-6 py-5";
const compactCellClassName = "px-3 py-1";
const emptyClassName = "h-32 px-6 text-center";

const disabledLinkClassName = "pointer-events-none opacity-50";

export const DataTable = <T,>({
  data,
  columns,
  emptyMessage = "No results.",
  compact = false,
  pageSize,
  getRowId,
}: DataTableProps<T>) => {
  const table = useReactTable({
    data: data || [],
    columns,
    getRowId,
    getCoreRowModel: getCoreRowModel(),
    ...(pageSize && {
      getPaginationRowModel: getPaginationRowModel(),
      autoResetPageIndex: false,
      initialState: { pagination: { pageIndex: 0, pageSize } },
    }),
  });

  const tableElement = (
    <Table>
      <TableHeader>
        {table.getHeaderGroups().map((headerGroup) => (
          <TableRow key={headerGroup.id} className={headRowClassName}>
            {headerGroup.headers.map((header) => {
              return (
                <TableHead
                  key={header.id}
                  className={compact ? compactHeadClassName : headClassName}
                >
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
                <TableCell
                  key={cell.id}
                  className={compact ? compactCellClassName : cellClassName}
                >
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

  if (!pageSize) {
    return tableElement;
  }

  return (
    <div className="flex flex-col gap-2">
      {tableElement}
      <Pagination className="justify-end">
        <PaginationContent>
          <PaginationItem>
            <PaginationPrevious
              onClick={() => table.previousPage()}
              className={
                table.getCanPreviousPage() ? undefined : disabledLinkClassName
              }
            />
          </PaginationItem>
          <PaginationItem className="px-2 text-sm text-muted-foreground">
            Page {table.getState().pagination.pageIndex + 1} of{" "}
            {Math.max(table.getPageCount(), 1)}
          </PaginationItem>
          <PaginationItem>
            <PaginationNext
              onClick={() => table.nextPage()}
              className={
                table.getCanNextPage() ? undefined : disabledLinkClassName
              }
            />
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    </div>
  );
};
