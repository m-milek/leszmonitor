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

  const onMessage = useCallback((event: MessageEvent) => {
    const data = JSON.parse(event.data);
    console.log("Received WebSocket message:", data);
  }, []);

  const { readyState } = useWebSocket(WEBSOCKET_ENDPOINT, {
    share: true,
    onMessage,
    filter: () => false,
    shouldReconnect: () => true,
    reconnectAttempts: 10,
    reconnectInterval: (attempt) =>
      Math.min(Math.pow(2, attempt) * 1000, 10000),
    heartbeat: {
      message: "ping",
      returnMessage: "pong",
      timeout: 5000,
      interval: 10000,
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
    setConnectionStatus(readyState);
  }, [readyState, setConnectionStatus]);

  return <>{children}</>;
}
