import { type ColumnDef } from "@tanstack/react-table";
import { NoData } from "@/components/common/NoData";
import type { StatusPage } from "@/features/status-pages/types";
import { DataTable, dataTableFeatures } from "@/components/common/DataTable";
import { ShortId } from "@/components/common/ShortId";
import { formatDate } from "@/lib/utils";

export interface StatusPagesTableProps {
  statusPages: StatusPage[];
}

const columns: ColumnDef<typeof dataTableFeatures, StatusPage>[] = [
  {
    accessorKey: "title",
    header: "Title",
  },
  {
    accessorKey: "slug",
    header: "Slug",
    cell: ({ row }) => (
      <code className="text-muted-foreground">{row.original.slug}</code>
    ),
  },
  {
    accessorKey: "description",
    header: "Description",
    cell: ({ row }) =>
      row.original.description ? (
        <span>{row.original.description}</span>
      ) : (
        <NoData label="No description" />
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
];

export const StatusPagesTable = ({ statusPages }: StatusPagesTableProps) => {
  return (
    <DataTable
      data={statusPages}
      columns={columns}
      emptyMessage="No status pages yet."
    />
  );
};
