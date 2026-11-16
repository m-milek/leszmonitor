import { MetadataApi } from "@/features/instance/metadata-api";
import { mapUserRoleToDisplayName, type User } from "@/features/users/types";
import {
  LucideChevronsUpDown,
  LucideLogOut,
  LucideUser,
  LucideUserCog,
} from "lucide-react";
import { useNavigate, useRouter } from "@tanstack/react-router";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItemIcon,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { Skeleton } from "@/components/ui/skeleton";
import { Flex } from "@/components/common/Flex";
import { Metadata } from "@/components/layout/Metadata";
import { useQuery } from "@tanstack/react-query";
import { clearToken } from "@/features/auth/lib/token";
import { Initial } from "@/features/users/components/Initial";

export interface AppSidebarFooterProps {
  user: User;
}

export const AppSidebarFooter = ({ user }: AppSidebarFooterProps) => {
  const router = useRouter();
  const navigate = useNavigate();

  const logOut = async () => {
    await clearToken();
    router.invalidate();
  };

  const { data: metadata } = useQuery({
    queryFn: async () => MetadataApi.get(),
    queryKey: ["metadata"],
  });

  return (
    <Flex direction="column">
      <SidebarMenu>
        <SidebarMenuItem>
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <SidebarMenuButton
                  size="lg"
                  aria-label="User menu"
                  className="data-open:bg-sidebar-accent"
                />
              }
            >
              <Initial text={user.username} size="sm" className="shrink-0" />
              <div className="grid flex-1 text-left leading-tight">
                <span className="truncate font-medium">{user.username}</span>
                <span className="truncate text-xs text-sidebar-foreground/70">
                  {mapUserRoleToDisplayName[user.role]}
                </span>
              </div>
              <LucideChevronsUpDown className="ml-auto" />
            </DropdownMenuTrigger>
            <DropdownMenuContent side="top" align="end" className="min-w-56">
              <DropdownMenuItemIcon
                icon={LucideUser}
                onClick={() =>
                  navigate({
                    to: "/user/$username",
                    params: { username: user.username },
                  })
                }
              >
                Profile
              </DropdownMenuItemIcon>
              <DropdownMenuItemIcon
                icon={LucideUserCog}
                onClick={() =>
                  navigate({
                    to: "/user/$username/settings",
                    params: { username: user.username },
                  })
                }
              >
                Settings
              </DropdownMenuItemIcon>
              <DropdownMenuItemIcon
                icon={LucideLogOut}
                variant="destructive"
                onClick={logOut}
              >
                Log out
              </DropdownMenuItemIcon>
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarMenuItem>
      </SidebarMenu>
      <div className="group-data-[collapsible=icon]:hidden">
        <Metadata data={metadata} />
      </div>
    </Flex>
  );
};

export const AppSidebarFooterSkeleton = () => (
  <div className="flex h-12 items-center gap-2 p-2 group-data-[collapsible=icon]:p-0">
    <Skeleton className="size-8 shrink-0 rounded-full" />
    <div className="grid flex-1 gap-1 group-data-[collapsible=icon]:hidden">
      <Skeleton className="h-3.5 w-24" />
      <Skeleton className="h-3 w-14" />
    </div>
  </div>
);
