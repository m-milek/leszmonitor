import type { User } from "@/lib/types.ts";
import { LucideEllipsisVertical } from "lucide-react";
import { Button } from "@/components/ui/button.tsx";

export interface AppSidebarFooterProps {
  user: User;
}

export const AppSidebarFooter = ({ user }: AppSidebarFooterProps) => {
  const firstLetter = user.username[0].toUpperCase();

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
        <Button variant="ghost">
          <LucideEllipsisVertical />
        </Button>
      </div>
    </div>
  );
};
