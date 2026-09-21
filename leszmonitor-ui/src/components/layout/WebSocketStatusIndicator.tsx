import { useAppStore } from "@/app/store";
import { ReadyState } from "react-use-websocket";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { StatusDot, type Status } from "@/components/common/StatusDot";
import type { WebSocketStatus } from "@/app/providers/websocket-status";

const connectionStatusLabel = {
  [ReadyState.CONNECTING]: "Connecting",
  [ReadyState.OPEN]: "OK",
  [ReadyState.CLOSING]: "Closing",
  [ReadyState.CLOSED]: "Closed",
  [ReadyState.UNINSTANTIATED]: "Uninstantiated",
};

const connectionStatusDot: Record<ReadyState, Status> = {
  [ReadyState.CONNECTING]: "pending",
  [ReadyState.OPEN]: "up",
  [ReadyState.CLOSING]: "pending",
  [ReadyState.CLOSED]: "down",
  [ReadyState.UNINSTANTIATED]: "unknown",
};

interface WebSocketStatusDisplayConfig {
  label: string;
  status: Status;
}

const displayWebSocketStatus = (
  wsStatus: WebSocketStatus,
): WebSocketStatusDisplayConfig => {
  const label = connectionStatusLabel[wsStatus.status];
  if (wsStatus.status === ReadyState.CLOSED) {
    return { label, status: "down" };
  }
  if (!wsStatus.isAuthenticated) {
    return { label: `${label} (Unauthenticated)`, status: "pending" };
  }
  return {
    label,
    status: connectionStatusDot[wsStatus.status] ?? "unknown",
  };
};

export const WebSocketStatusIndicator = () => {
  const { webSocketConnectionStatus: wsStatus } = useAppStore();
  const { label, status } = displayWebSocketStatus(wsStatus);
  return (
    <Tooltip>
      <TooltipTrigger delay={500}>
        <StatusDot status={status} />
      </TooltipTrigger>
      <TooltipContent side="top">Connection Status: {label}</TooltipContent>
    </Tooltip>
  );
};
