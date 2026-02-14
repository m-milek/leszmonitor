import { type ColumnDef, getCoreRowModel } from "@tanstack/table-core";
import type { TeamMember } from "@/lib/types.ts";
import { StyledLink } from "@/components/leszmonitor/StyledLink.tsx";
import { flexRender, useReactTable } from "@tanstack/react-table";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx";
import { formatDate } from "@/lib/utils.ts";

const columns: ColumnDef<TeamMember>[] = [
  {
    accessorKey: "username",
    header: "User",
    cell: ({ row }) => {
      return (
        <StyledLink
          to="/user/$username"
          params={{ username: row.original.username }}
        >
          {row.original.username}
        </StyledLink>
      );
    },
  },
  {
    accessorKey: "role",
    header: "Role",
  },
  {
    accessorKey: "createdAt",
    header: "Joined at",
    cell: ({ row }) => {
      return formatDate(row.original.createdAt);
    },
  },
  {
    accessorKey: "updatedAt",
    header: "Last updated",
    cell: ({ row }) => {
      return formatDate(row.original.updatedAt);
    },
  },
];

export interface TeamMembersTableProps {
  members: TeamMember[];
}

export const TeamMembersTable = ({ members }: TeamMembersTableProps) => {
  const table = useReactTable({
    data: members,
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
