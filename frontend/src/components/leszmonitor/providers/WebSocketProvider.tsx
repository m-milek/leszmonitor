import { type ReactNode, useCallback, useEffect, useState } from "react";
import useWebSocket from "react-use-websocket";
import { WEBSOCKET_ENDPOINT } from "@/lib/data/webSocket.ts";
import { webSocketConnectionStatusAtom } from "@/lib/atoms.ts";
import { useSetAtom } from "jotai";
import { getLoginToken } from "@/lib/utils.ts";
import { toast } from "sonner";

type WebSocketProviderProps = {
  children: ReactNode;
};

export function WebSocketProvider({ children }: WebSocketProviderProps) {
  const setConnectionStatus = useSetAtom(webSocketConnectionStatusAtom);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

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
      console.log("Received WebSocket message:", data);
      toast.info("Received WebSocket message");
    } catch {
      console.log("Received WebSocket message:", event.data);
    }
  }, []);

  const { readyState, sendMessage, getWebSocket } = useWebSocket(
    WEBSOCKET_ENDPOINT,
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
