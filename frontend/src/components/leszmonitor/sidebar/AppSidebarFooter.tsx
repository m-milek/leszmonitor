import type { User } from "@/lib/types.ts";
import { LucideEllipsisVertical, LucideLogOut } from "lucide-react";
import { Button } from "@/components/ui/button.tsx";
import { useRouter } from "@tanstack/react-router";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { Metadata } from "@/components/leszmonitor/sidebar/Metadata.tsx";
import { useQuery } from "@tanstack/react-query";
import { getMetadata } from "@/lib/data/metadata-api.ts";

export interface AppSidebarFooterProps {
  user: User;
}

export const AppSidebarFooter = ({ user }: AppSidebarFooterProps) => {
  const firstLetter = user.username[0].toUpperCase();

  const router = useRouter();

  const logOut = async () => {
    await cookieStore.delete("LOGIN_TOKEN");
    router.invalidate();
  };

  const { data: metadata } = useQuery({
    queryFn: async () => getMetadata(),
    queryKey: ["metadata"],
  });

  return (
    <Flex direction="column">
      <div className="flex items-center p-4">
        <div className="flex h-9 w-9 items-center justify-center rounded-full bg-primary">
          <span className="text-sm font-medium text-white">{firstLetter}</span>
        </div>
        <div className="flex flex-1 items-center justify-between">
          <div className="ml-2">
            <p className="font-medium">{user.username}</p>
            <p className="text-sm">Logged in</p>
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost">
                <LucideEllipsisVertical />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem className="text-destructive" onSelect={logOut}>
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
