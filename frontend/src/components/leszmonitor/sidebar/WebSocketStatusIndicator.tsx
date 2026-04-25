import { useAtomValue } from "jotai";
import { webSocketConnectionStatusAtom } from "@/lib/atoms.ts";
import { ReadyState } from "react-use-websocket";
import { cn } from "@/lib/utils.ts";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip.tsx";

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

export const WebSocketStatusIndicator = () => {
  const wsStatus = useAtomValue(webSocketConnectionStatusAtom);
  return (
    <Tooltip delayDuration={500}>
      <TooltipTrigger>
        <div
          className={cn(
            "h-4 w-4 rounded-full",
            connectionStatusColor[wsStatus],
          )}
        />
      </TooltipTrigger>
      <TooltipContent side="top">
        WebSocket Status: {connectionStatusLabel[wsStatus]}
      </TooltipContent>
    </Tooltip>
  );
};
