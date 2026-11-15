import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  type ColumnDef,
  createPaginatedRowModel,
  createSortedRowModel,
  flexRender,
  type RowData,
  rowPaginationFeature,
  rowSortingFeature,
  tableFeatures,
  useTable,
} from "@tanstack/react-table";
import { ArrowDown, ArrowUp, ArrowUpDown } from "lucide-react";
import type { ReactNode } from "react";
import { cn } from "cn";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

export const dataTableFeatures = tableFeatures({
  rowPaginationFeature,
  rowSortingFeature,
  paginatedRowModel: createPaginatedRowModel(),
  sortedRowModel: createSortedRowModel(),
});

export interface DataTableProps<T extends RowData> {
  data: T[];
  columns: ColumnDef<typeof dataTableFeatures, T>[];
  emptyMessage?: ReactNode;
  compact?: boolean;
  pageSize?: number;
  getRowId?: (row: T) => string;
  sortable?: boolean;
  onRowClick?: (row: T) => void;
}

const headRowClassName = "hover:bg-transparent";
const headClassName = "h-12 px-6";
const compactHeadClassName = "h-9 px-3";
const bodyRowClassName = "transition-colors hover:bg-muted/40";
const cellClassName = "px-6 py-5";
const compactCellClassName = "px-3 py-1";
const emptyClassName = "h-32 px-6 text-center";

const disabledLinkClassName = "pointer-events-none opacity-50";

const SortIcon = ({ direction }: { direction: false | "asc" | "desc" }) => {
  if (direction === "asc") return <ArrowUp className="size-3.5" />;
  if (direction === "desc") return <ArrowDown className="size-3.5" />;
  return <ArrowUpDown className="size-3.5 opacity-40" />;
};

export const DataTable = <T extends RowData>({
  data,
  columns,
  emptyMessage = "No results.",
  compact = false,
  pageSize,
  getRowId,
  sortable = false,
  onRowClick,
}: DataTableProps<T>) => {
  const table = useTable({
    features: dataTableFeatures,
    data: data || [],
    columns,
    getRowId,
    enableSorting: sortable,
    enableMultiSort: false,
    manualPagination: !pageSize,
    ...(pageSize && {
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
                  {header.isPlaceholder ? null : header.column.getCanSort() ? (
                    <button
                      type="button"
                      className="inline-flex cursor-pointer items-center gap-1 hover:text-foreground"
                      onClick={header.column.getToggleSortingHandler()}
                    >
                      {flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                      <SortIcon direction={header.column.getIsSorted()} />
                    </button>
                  ) : (
                    flexRender(
                      header.column.columnDef.header,
                      header.getContext(),
                    )
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
              className={cn(bodyRowClassName, onRowClick && "cursor-pointer")}
              onClick={onRowClick && (() => onRowClick(row.original))}
            >
              {row.getAllCells().map((cell) => (
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
            Page {table.state.pagination.pageIndex + 1} of{" "}
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
