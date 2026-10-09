import { type ColumnDef } from "@tanstack/table-core";
import { NoData } from "@/components/common/NoData";
import type { Tag as TagModel } from "@/features/tags/types";
import { DataTable } from "@/components/common/DataTable";
import { Tag } from "@/features/tags/components/Tag";
import { ShortId } from "@/components/common/ShortId";
import { formatDate } from "@/lib/utils";
import { normalizeHexColor } from "@/features/tags/lib/colors";
import { DeleteTagDialog } from "@/features/tags/components/DeleteTagDialog";
import { Flex } from "@/components/common/Flex";
import { Badge } from "@/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { LucideLock } from "lucide-react";

export interface TagsTableProps {
  tags: TagModel[];
}

const columns: ColumnDef<TagModel>[] = [
  {
    accessorKey: "name",
    header: "Tag",
    cell: ({ row }) => (
      <Flex direction="row" className="items-center gap-2">
        <Tag tag={row.original} />
        {row.original.source === "config" && (
          <Tooltip>
            <TooltipTrigger>
              <Badge variant="ghost" className="h-7 px-2 [&>svg]:size-5!">
                <LucideLock />
              </Badge>
            </TooltipTrigger>
            <TooltipContent>Tag defined in config file</TooltipContent>
          </Tooltip>
        )}
      </Flex>
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
    <DataTable data={tags} columns={columns} emptyMessage="No tags yet." />
  );
};
