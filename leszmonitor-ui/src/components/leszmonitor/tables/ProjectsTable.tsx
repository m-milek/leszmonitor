import type { Project } from "@/lib/types";
import { formatDate } from "@/lib/utils.ts";
import { StyledLink } from "../StyledLink";
import { type ColumnDef, type Row } from "@tanstack/table-core";
import { Button } from "@/components/ui/button.tsx";
import { LucideTrash2 } from "lucide-react";
import { GenericTable } from "@/components/leszmonitor/tables/GenericTable.tsx";

export interface ProjectsTableProps {
  projects: Project[];
  onProjectDeleted?: (projectId: string) => Promise<void>;
}

export const ProjectsTable = ({
  projects,
  onProjectDeleted,
}: ProjectsTableProps) => {
  const columns: ColumnDef<Project>[] = [
    {
      accessorKey: "name",
      header: "Name",
      cell: ({ row }) => {
        const name = row.original.name;
        return (
          <StyledLink
            to="/projects/$projectId"
            params={{ projectId: row.original.slug }}
          >
            {name}
          </StyledLink>
        );
      },
    },
    {
      accessorKey: "slug",
      header: "Slug",
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
    ...(onProjectDeleted
      ? [
          {
            header: "Actions",
            cell: ({ row }: { row: Row<Project> }) => {
              const projectId = row.original.slug;
              return (
                <Button
                  variant="ghost"
                  onClick={() => onProjectDeleted(projectId)}
                  className="text-destructive"
                >
                  <LucideTrash2 />
                </Button>
              );
            },
          },
        ]
      : []),
  ];

  return <GenericTable data={projects} columns={columns} />;
};
