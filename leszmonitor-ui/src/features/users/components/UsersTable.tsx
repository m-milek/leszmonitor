import {
  UserRole,
  mapUserRoleToDisplayName,
  type User,
} from "@/features/users/types";
import { formatDate } from "@/lib/utils";
import { StyledLink } from "@/components/common/StyledLink";
import { type ColumnDef } from "@tanstack/table-core";
import { DataTable } from "@/components/common/DataTable";
import { MoreVertical, Trash2 } from "lucide-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { removeUser, updateUserRole } from "@/features/users/users-api";
import { LMSelect } from "@/components/form/LMSelect";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItemIcon,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

const roleSelectItems = Object.values(UserRole).map((role) => ({
  value: role,
  label: mapUserRoleToDisplayName[role],
}));

const RoleCell = ({ user }: { user: User }) => {
  const queryClient = useQueryClient();
  const roleMutation = useMutation({
    mutationFn: (role: UserRole) => updateUserRole(user.username, role),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });

  return (
    <LMSelect
      id={`role-${user.username}`}
      name={`role-${user.username}`}
      value={user.role}
      onValueChange={(value) => roleMutation.mutate(value as UserRole)}
      items={roleSelectItems}
    />
  );
};

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
      accessorKey: "role",
      header: "Role",
      cell: ({ row }) => <RoleCell user={row.original} />,
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

  return <DataTable data={users} columns={columns} />;
};
