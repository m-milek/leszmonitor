import { MetadataApi } from "@/features/instance/metadata-api";
import type { User } from "@/features/users/types";
import { LucideEllipsisVertical, LucideLogOut } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useRouter } from "@tanstack/react-router";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
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
      <div className="flex items-center">
        <Initial text={user.username} size="sm" />
        <div className="flex flex-1 items-center justify-between">
          <div className="ml-2">
            <p className="font-medium">{user.username}</p>
            <p className="text-sm">Logged in</p>
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger render={<Button variant="ghost" />}>
              <LucideEllipsisVertical />
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem className="text-destructive" onClick={logOut}>
                <div className="flex items-center w-full justify-between">
                  <span>Log out</span>
                  <LucideLogOut className="text-destructive" />
                </div>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
      <Metadata data={metadata} />
    </Flex>
  );
};
