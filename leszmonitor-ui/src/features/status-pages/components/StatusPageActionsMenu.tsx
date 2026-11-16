import { useState } from "react";
import {
  FileOutputIcon,
  LucidePencil,
  MoreVertical,
  Trash2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItemIcon,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { EditStatusPageDialog } from "@/features/status-pages/components/EditStatusPageDialog";
import { DeleteStatusPageDialog } from "@/features/status-pages/components/DeleteStatusPageDialog";
import { ExportStatusPageDialog } from "@/features/status-pages/components/ExportStatusPageDialog";
import type { StatusPage } from "@/features/status-pages/types";

export interface StatusPageActionsMenuProps {
  statusPage: StatusPage;
}

export function StatusPageActionsMenu({
  statusPage,
}: StatusPageActionsMenuProps) {
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [isDeleteOpen, setIsDeleteOpen] = useState(false);
  const [isExportOpen, setIsExportOpen] = useState(false);
  const isConfigBased = statusPage.source === "config";

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
            disabled={isConfigBased}
            onClick={() => setIsEditOpen(true)}
          >
            Edit
          </DropdownMenuItemIcon>
          <DropdownMenuItemIcon
            icon={FileOutputIcon}
            onClick={() => setIsExportOpen(true)}
          >
            Export
          </DropdownMenuItemIcon>
          <DropdownMenuItemIcon
            icon={Trash2}
            variant="destructive"
            disabled={isConfigBased}
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
      <ExportStatusPageDialog
        statusPage={statusPage}
        open={isExportOpen}
        onOpenChange={setIsExportOpen}
      />
      <DeleteStatusPageDialog
        statusPage={statusPage}
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
      />
    </>
  );
}
