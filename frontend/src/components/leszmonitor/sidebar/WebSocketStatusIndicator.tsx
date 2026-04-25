import { useAtomValue } from "jotai";
import { webSocketConnectionStatusAtom } from "@/lib/atoms.ts";
import { ReadyState } from "react-use-websocket";
import { cn } from "@/lib/utils.ts";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip.tsx";
import type { WebSocketStatus } from "@/lib/data/webSocket.ts";

const connectionStatusLabel = {
  [ReadyState.CONNECTING]: "Connecting",
  [ReadyState.OPEN]: "Open",
  [ReadyState.CLOSING]: "Closing",
  [ReadyState.CLOSED]: "Closed",
  [ReadyState.UNINSTANTIATED]: "Uninstantiated",
};

const connectionStatusColor = {
  [ReadyState.CONNECTING]: "bg-yellow-500",
  [ReadyState.OPEN]: "bg-green-500",
  [ReadyState.CLOSING]: "bg-orange-500",
  [ReadyState.CLOSED]: "bg-red-700",
  [ReadyState.UNINSTANTIATED]: "bg-gray-500",
};

interface WebSocketStatusDisplayConfig {
  label: string;
  colorClass: string;
}

const displayWebSocketStatus = (
  wsStatus: WebSocketStatus,
): WebSocketStatusDisplayConfig => {
  const label = connectionStatusLabel[wsStatus.status];
  if (!wsStatus.isAuthenticated) {
    return { label: `${label} (Unauthenticated)`, colorClass: "bg-yellow-500" };
  }
  return {
    label,
    colorClass: connectionStatusColor[wsStatus.status] || "bg-gray-500",
  };
};

export const WebSocketStatusIndicator = () => {
  const wsStatus = useAtomValue(webSocketConnectionStatusAtom);
  const { label, colorClass } = displayWebSocketStatus(wsStatus);
  return (
    <Tooltip delayDuration={500}>
      <TooltipTrigger>
        <div className={cn("h-3 w-3 rounded-full", colorClass)} />
      </TooltipTrigger>
      <TooltipContent side="top">WebSocket Status: {label}</TooltipContent>
    </Tooltip>
  );
};
