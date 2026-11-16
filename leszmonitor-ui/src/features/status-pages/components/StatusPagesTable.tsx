import { useMemo } from "react";
import { Link } from "@tanstack/react-router";
import { type ColumnDef } from "@tanstack/react-table";
import { LucideExternalLink } from "lucide-react";
import { DataTable, dataTableFeatures } from "@/components/common/DataTable";
import { NoData } from "@/components/common/NoData";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { Tag } from "@/features/tags/components/Tag";
import { StatusPageActionsMenu } from "@/features/status-pages/components/StatusPageActionsMenu";
import type { Tag as TagModel } from "@/features/tags/types";
import type { StatusPage } from "@/features/status-pages/types";

export interface StatusPagesTableProps {
  statusPages: StatusPage[];
  tags: TagModel[];
}

export const StatusPagesTable = ({
  statusPages,
  tags,
}: StatusPagesTableProps) => {
  const columns = useMemo<ColumnDef<typeof dataTableFeatures, StatusPage>[]>(
    () => [
      {
        accessorKey: "title",
        header: "Title",
        cell: ({ row }) => (
          <span className="font-medium">{row.original.title}</span>
        ),
      },
      {
        id: "tags",
        header: "Tags",
        cell: ({ row }) => {
          const pageTags = tags.filter((tag) =>
            row.original.tagIds.includes(tag.id),
          );
          return pageTags.length > 0 ? (
            <div className="flex flex-wrap gap-1">
              {pageTags.map((tag) => (
                <Tag key={tag.id} tag={tag} />
              ))}
            </div>
          ) : (
            <NoData label="No tags" />
          );
        },
      },
      {
        accessorKey: "monitorCount",
        header: "Monitors",
        cell: ({ row }) => (
          <span className="text-sm">{row.original.monitorCount}</span>
        ),
      },
      {
        accessorKey: "enabled",
        header: "Status",
        cell: ({ row }) => (
          <span className="text-sm">
            {row.original.enabled ? "Enabled" : "Disabled"}
          </span>
        ),
      },
      {
        accessorKey: "requireAuth",
        header: "Access",
        cell: ({ row }) => (
          <span className="text-sm">
            {row.original.requireAuth ? "Login required" : "Public"}
          </span>
        ),
      },
      {
        accessorKey: "groupMonitorsByTag",
        header: "Grouped by tag",
        cell: ({ row }) => (
          <span className="text-sm">
            {row.original.groupMonitorsByTag ? "Yes" : "No"}
          </span>
        ),
      },
      {
        id: "actions",
        header: "",
        cell: ({ row }) => (
          <div className="flex justify-end">
            <Tooltip>
              <TooltipTrigger
                render={
                  <Button
                    variant="ghost"
                    size="icon"
                    aria-label="Open status page"
                    nativeButton={false}
                    render={
                      <Link
                        to="/status/$slug"
                        params={{ slug: row.original.slug }}
                        target="_blank"
                      />
                    }
                  />
                }
              >
                <LucideExternalLink />
              </TooltipTrigger>
              <TooltipContent>Open status page</TooltipContent>
            </Tooltip>
            <StatusPageActionsMenu statusPage={row.original} />
          </div>
        ),
      },
    ],
    [tags],
  );

  return (
    <DataTable
      data={statusPages}
      columns={columns}
      emptyMessage="No status pages yet."
    />
  );
};
