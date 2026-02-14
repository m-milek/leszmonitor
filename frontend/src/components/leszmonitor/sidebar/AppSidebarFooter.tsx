import type { User } from "@/lib/types.ts";
import { LucideEllipsisVertical } from "lucide-react";
import { Button } from "@/components/ui/button.tsx";
import { useRouter } from "@tanstack/react-router";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";

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

  return (
    <div className="flex items-center m-2">
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
            <DropdownMenuLabel>
              <button onClick={logOut}>Log out</button>
            </DropdownMenuLabel>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );
};
