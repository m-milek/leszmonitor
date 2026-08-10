import type { User } from "@/lib/types.ts";
import { formatDate } from "@/lib/utils.ts";
import { StyledLink } from "../StyledLink";
import { type ColumnDef } from "@tanstack/table-core";
import { GenericTable } from "@/components/leszmonitor/tables/GenericTable.tsx";
import { MoreVertical, Pencil, Trash2 } from "lucide-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { removeUser } from "@/lib/data/userData.ts";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItemIcon,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

const ActionsCell = ({ user }: { user: User }) => {
  const queryClient = useQueryClient();
  const removeMutation = useMutation({
    mutationFn: () => removeUser(user.username),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button className="flex items-center justify-center p-2 outline-none rounded-md hover:bg-accent hover:text-accent-foreground">
          <MoreVertical className="w-4 h-4" />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItemIcon
          icon={Pencil}
          onClick={() => console.log("Edit user", user.username)}
        >
          Edit
        </DropdownMenuItemIcon>
        <DropdownMenuItemIcon
          icon={Trash2}
          className="text-destructive focus:text-destructive"
          onClick={() => removeMutation.mutate()}
        >
          Delete
        </DropdownMenuItemIcon>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export interface UsersTableProps {
  users: User[];
}

export const UsersTable = ({ users }: UsersTableProps) => {
  const columns: ColumnDef<User>[] = [
    {
      accessorKey: "username",
      header: "Username",
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
      accessorKey: "id",
      header: "ID",
    },
    {
      accessorKey: "createdAt",
      header: "Joined At",
      cell: ({ row }) => {
        return formatDate(row.original.createdAt);
      },
    },
    {
      accessorKey: "updatedAt",
      header: "Last Updated",
      cell: ({ row }) => {
        return formatDate(row.original.updatedAt);
      },
    },
    {
      accessorKey: "",
      id: "contextMenu",
      cell: ({ row }) => {
        return <ActionsCell user={row.original} />;
      },
    },
  ];

  return <GenericTable data={users} columns={columns} />;
};
