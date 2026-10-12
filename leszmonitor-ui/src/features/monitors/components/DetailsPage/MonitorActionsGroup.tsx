import { ButtonGroup } from "@/components/ui/button-group";
import { Button } from "@/components/ui/button";
import {
  EllipsisIcon,
  FileOutputIcon,
  LucideCirclePlay,
  PauseIcon,
  PencilIcon,
  PlayIcon,
} from "lucide-react";
import { useState } from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItemIcon,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { DeleteMonitorDialog } from "@/features/monitors/components/DeleteMonitorDialog.tsx";
import type { Monitor } from "@/features/monitors/types.ts";
import { useNavigate } from "@tanstack/react-router";
import { ExportMonitorDialog } from "@/features/monitors/components/DetailsPage/ExportMonitorDialog";

export interface MonitorActionsGroupProps {
  monitor: Monitor;
  handleToggleMonitorState: () => void;
  handleEditMonitor: () => void;
  handleManuallyRunMonitor: () => void;
  isPaused: boolean;
}

export const MonitorActionsGroup = ({
  monitor,
  handleToggleMonitorState,
  handleEditMonitor,
  handleManuallyRunMonitor,
  isPaused,
}: MonitorActionsGroupProps) => {
  const navigate = useNavigate();
  const [isExportOpen, setIsExportOpen] = useState(false);
  const isConfigBased = monitor.source === "config";
  return (
    <>
      <ButtonGroup>
        <Button
          variant="outline"
          size="lg"
          onClick={handleToggleMonitorState}
          disabled={isConfigBased}
        >
          {isPaused ? (
            <>
              <PlayIcon /> Resume
            </>
          ) : (
            <>
              <PauseIcon /> Pause
            </>
          )}
        </Button>
        <Button
          variant="outline"
          size="lg"
          onClick={handleEditMonitor}
          disabled={isConfigBased}
        >
          <PencilIcon />
          <span>Edit</span>
        </Button>
        {monitor.type !== "push" && (
          <Button
            variant="outline"
            size="lg"
            onClick={handleManuallyRunMonitor}
          >
            <LucideCirclePlay />
            <span>Run Now</span>
          </Button>
        )}
        <DeleteMonitorDialog
          monitor={monitor}
          onDeleted={() => navigate({ to: "/monitors" })}
        />
        <DropdownMenu>
          <DropdownMenuTrigger
            render={
              <Button
                variant="outline"
                size="icon-lg"
                aria-label="More actions"
              />
            }
          >
            <EllipsisIcon />
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItemIcon
              icon={FileOutputIcon}
              onClick={() => setIsExportOpen(true)}
            >
              Export
            </DropdownMenuItemIcon>
          </DropdownMenuContent>
        </DropdownMenu>
      </ButtonGroup>
      <ExportMonitorDialog
        monitor={monitor}
        open={isExportOpen}
        onOpenChange={setIsExportOpen}
      />
    </>
  );
};
export default MonitorActionsGroup;
