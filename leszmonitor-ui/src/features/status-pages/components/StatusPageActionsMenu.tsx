import { useState } from "react";
import { LucidePencil, MoreVertical, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItemIcon,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { EditStatusPageDialog } from "@/features/status-pages/components/EditStatusPageDialog";
import { DeleteStatusPageDialog } from "@/features/status-pages/components/DeleteStatusPageDialog";
import type { StatusPage } from "@/features/status-pages/types";

export interface StatusPageActionsMenuProps {
  statusPage: StatusPage;
}

export function StatusPageActionsMenu({
  statusPage,
}: StatusPageActionsMenuProps) {
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger
          render={
            <Button variant="ghost" size="icon" aria-label="More actions" />
          }
        >
          <MoreVertical />
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItemIcon
            icon={LucidePencil}
            onClick={() => setIsEditOpen(true)}
          >
            Edit
          </DropdownMenuItemIcon>
          <DropdownMenuItemIcon
            icon={Trash2}
            variant="destructive"
            onClick={() => setIsDeleteOpen(true)}
          >
            Delete
          </DropdownMenuItemIcon>
        </DropdownMenuContent>
      </DropdownMenu>
      <EditStatusPageDialog
        statusPage={statusPage}
        open={isEditOpen}
        onOpenChange={setIsEditOpen}
      />
      <DeleteStatusPageDialog
        statusPage={statusPage}
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
      />
    </>
  );
}
