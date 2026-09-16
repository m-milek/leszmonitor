import { type ColumnDef } from "@tanstack/table-core";
import { Minus } from "lucide-react";
import type { Tag as TagModel } from "@/features/tags/types";
import { DataTable } from "@/components/common/DataTable";
import { Tag } from "@/features/tags/components/Tag";
import { ShortId } from "@/components/common/ShortId";
import { formatDate } from "@/lib/utils";
import { normalizeHexColor } from "@/features/tags/lib/colors";
import { DeleteTagDialog } from "@/features/tags/components/DeleteTagDialog";

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
  return (
    <DataTable
      data={tags}
      columns={columns}
      wrapperClassName="rounded-md border border-border bg-card"
      headRowClassName="hover:bg-transparent"
      headClassName="h-12 px-6"
      bodyRowClassName="transition-colors hover:bg-muted/40"
      cellClassName="px-6 py-5"
      emptyMessage="No tags yet."
      emptyClassName="h-32 px-6 text-center"
    />
  );
};
