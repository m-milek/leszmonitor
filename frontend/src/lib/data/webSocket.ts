import { BACKEND_WS_URL } from "@/lib/consts.ts";
import type { ReadyState } from "react-use-websocket";

export const WEBSOCKET_ENDPOINT = `${BACKEND_WS_URL}/ws`;

export interface WebSocketStatus {
  status: ReadyState;
  isAuthenticated: boolean;
}
