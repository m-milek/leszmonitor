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
import { ExportYamlDialog } from "@/components/common/ExportYamlDialog";
import { StatusPagesApi } from "@/features/status-pages/status-pages-api";
import { QUERY_KEYS } from "@/lib/consts";
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
      <ExportYamlDialog
        title={`Export ${statusPage.title}`}
        queryKey={[QUERY_KEYS.STATUS_PAGES, statusPage.id, "export"]}
        queryFn={() => StatusPagesApi.exportYaml(statusPage.id)}
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
