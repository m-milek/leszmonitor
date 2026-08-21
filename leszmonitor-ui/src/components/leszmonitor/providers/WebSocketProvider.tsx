import { type ReactNode, useCallback, useEffect, useState } from "react";
import useWebSocket from "react-use-websocket";
import { useAppStore } from "@/lib/store.ts";
import { getLoginToken } from "@/lib/utils.ts";
import { toast } from "sonner";
import { useQueryClient } from "@tanstack/react-query";
import { isMonitorResultMessage } from "@/lib/types.ts";
import { QUERY_KEYS, SERVER_WS_URL } from "@/lib/consts.ts";

type WebSocketProviderProps = {
  children: ReactNode;
};

export function WebSocketProvider({
  children,
}: Readonly<WebSocketProviderProps>) {
  const { setWebSocketConnectionStatus: setConnectionStatus } = useAppStore();
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const queryClient = useQueryClient();

  const onMessage = useCallback((event: MessageEvent) => {
    if (event.data === "pong") {
      return;
    }

    try {
      const data = JSON.parse(event.data);
      if (data?.type === "auth" && data?.status === "ok") {
        setIsAuthenticated(true);
        return;
      }

      if (isMonitorResultMessage(data)) {
        if (data.response.status === "up") {
          toast.success(`Monitor ${data.monitorId} succeeded`);
        } else if (data.response.status === "down") {
          toast.error(`Monitor ${data.monitorId} failed`);
        } else {
          toast.info(
            `Monitor ${data.monitorId} status: ${data.response.status}`,
          );
        }
        queryClient.invalidateQueries({
          queryKey: [QUERY_KEYS.MONITOR_RESULTS, data.monitorId],
        });
      }
    } catch {
      console.log("Received WebSocket message:", event.data);
    }
  }, []);

  const { readyState, sendMessage, getWebSocket } = useWebSocket(
    SERVER_WS_URL,
    {
      share: true,
      onMessage,
      shouldReconnect: () => true,
      reconnectAttempts: 10,
      reconnectInterval: (attempt) =>
        Math.min(Math.pow(2, attempt) * 1000, 10000),
      heartbeat: {
        message: "ping",
        returnMessage: "pong",
        interval: 5000,
        timeout: 15000,
      },
      onOpen: () => {
        console.log("WebSocket connection opened");
        setIsAuthenticated(false);

        void (async () => {
          const token = await getLoginToken();
          if (!token) {
            const ws = getWebSocket();
            if (ws instanceof WebSocket) {
              ws.close(1008, "Missing auth token");
            }
            return;
          }

          sendMessage(
            JSON.stringify({
              type: "auth",
              token,
            }),
          );
        })();
      },
      onClose: () => {
        setIsAuthenticated(false);
        console.log("WebSocket connection closed");
      },
      onError: (event) => {
        setIsAuthenticated(false);
        console.error("WebSocket error:", event);
      },
    },
  );

  useEffect(() => {
    setConnectionStatus({
      status: readyState,
      isAuthenticated,
    });
  }, [readyState, isAuthenticated, setConnectionStatus]);

  return <>{children}</>;
}
