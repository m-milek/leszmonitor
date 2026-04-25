import { type ReactNode, useCallback, useEffect } from "react";
import useWebSocket from "react-use-websocket";
import { WEBSOCKET_ENDPOINT } from "@/lib/data/webSocket.ts";
import { webSocketConnectionStatusAtom } from "@/lib/atoms.ts";
import { useSetAtom } from "jotai";

type WebSocketProviderProps = {
  children: ReactNode;
};

export function WebSocketProvider({ children }: WebSocketProviderProps) {
  const setConnectionStatus = useSetAtom(webSocketConnectionStatusAtom);
  const isAuthenticated = false;

  const onMessage = useCallback((event: MessageEvent) => {
    if (event.data === "pong") {
      return;
    }

    try {
      const data = JSON.parse(event.data);
      console.log("Received WebSocket message:", data);
    } catch {
      console.log("Received WebSocket message:", event.data);
    }
  }, []);

  const { readyState } = useWebSocket(WEBSOCKET_ENDPOINT, {
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
    },
    onClose: () => {
      console.log("WebSocket connection closed");
    },
    onError: (event) => {
      console.error("WebSocket error:", event);
    },
  });

  useEffect(() => {
    setConnectionStatus({
      status: readyState,
      isAuthenticated,
    });
  }, [readyState, isAuthenticated, setConnectionStatus]);

  return <>{children}</>;
}
