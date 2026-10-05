import { ButtonGroup } from "@/components/ui/button-group";
import { Button } from "@/components/ui/button";
import {
  LucideCirclePlay,
  PauseIcon,
  PencilIcon,
  PlayIcon,
} from "lucide-react";
import { DeleteMonitorDialog } from "@/features/monitors/components/DeleteMonitorDialog.tsx";
import type { Monitor } from "@/features/monitors/types.ts";
import { useNavigate } from "@tanstack/react-router";

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

  return (
    <ButtonGroup>
      <Button variant="outline" size="lg" onClick={handleToggleMonitorState}>
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
      <Button variant="outline" size="lg" onClick={handleEditMonitor}>
        <PencilIcon />
        <span>Edit</span>
      </Button>
      {monitor.type !== "push" && (
        <Button variant="outline" size="lg" onClick={handleManuallyRunMonitor}>
          <LucideCirclePlay />
          <span>Run Now</span>
        </Button>
      )}
      <DeleteMonitorDialog
        monitor={monitor}
        onDeleted={() => navigate({ to: "/monitors" })}
      />
    </ButtonGroup>
  );
};
export default MonitorActionsGroup;
