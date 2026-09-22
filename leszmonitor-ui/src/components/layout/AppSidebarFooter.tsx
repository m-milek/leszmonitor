import { MetadataApi } from "@/features/instance/metadata-api";
import type { User } from "@/features/users/types";
import {
  LucideEllipsisVertical,
  LucideLogOut,
  LucideUserCog,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Link, useNavigate, useRouter } from "@tanstack/react-router";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItemIcon,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
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
      <div className="flex items-center gap-1">
        <Link
          to="/user/$username"
          params={{ username: user.username }}
          className="flex flex-1 items-center rounded-lg p-1 hover:bg-sidebar-accent"
        >
          <Initial text={user.username} size="sm" />
          <div className="ml-2">
            <p className="font-medium">{user.username}</p>
            <p className="text-sm">Logged in</p>
          </div>
        </Link>
        <DropdownMenu>
          <DropdownMenuTrigger
            render={<Button variant="ghost" aria-label="User menu" />}
          >
            <LucideEllipsisVertical />
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="rounded-[18px]">
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
      </div>
      <Metadata data={metadata} />
    </Flex>
  );
};
