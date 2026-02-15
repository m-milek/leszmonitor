import type { Group } from "@/lib/types";
import { formatDate } from "@/lib/utils.ts";
import { StyledLink } from "../StyledLink";
import { type ColumnDef, getCoreRowModel } from "@tanstack/table-core";
import { Button } from "@/components/ui/button.tsx";
import { flexRender, useReactTable } from "@tanstack/react-table";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx";
import { LucideTrash2 } from "lucide-react";

export interface GroupsTableProps {
  groups: Group[];
  teamId: string;
  onGroupDeleted: (groupId: string) => Promise<void>;
}

export const GroupsTable = ({
  groups,
  teamId,
  onGroupDeleted,
}: GroupsTableProps) => {
  const columns: ColumnDef<Group>[] = [
    {
      accessorKey: "name",
      header: "Name",
      cell: ({ row }) => {
        const name = row.original.name;
        return (
          <StyledLink
            to="/team/$teamId/groups/$groupId"
            params={{ teamId, groupId: row.original.displayId }}
          >
            {name}
          </StyledLink>
        );
      },
    },
    {
      accessorKey: "displayId",
      header: "Display ID",
    },
    {
      accessorKey: "createdAt",
      header: "Created At",
      cell: ({ row }) => {
        return formatDate(row.original.createdAt);
      },
    },
    {
      accessorKey: "updatedAt",
      header: "Updated At",
      cell: ({ row }) => {
        return formatDate(row.original.updatedAt);
      },
    },
    {
      accessorKey: "description",
      header: "Description",
    },
    {
      header: "Actions",
      cell: ({ row }) => {
        const groupId = row.original.displayId;
        return (
          <Button
            variant="ghost"
            onClick={() => onGroupDeleted(groupId)}
            className="text-destructive"
          >
            <LucideTrash2 />
          </Button>
        );
      },
    },
  ];

  const table = useReactTable({
    data: groups || [],
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <Table>
      <TableHeader>
        {table.getHeaderGroups().map((headerGroup) => (
          <TableRow key={headerGroup.id}>
            {headerGroup.headers.map((header) => {
              return (
                <TableHead key={header.id}>
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
            >
              {row.getVisibleCells().map((cell) => (
                <TableCell key={cell.id}>
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </TableCell>
              ))}
            </TableRow>
          ))
        ) : (
          <TableRow>
            <TableCell colSpan={columns.length} className="h-24 text-center">
              No results.
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );
};
