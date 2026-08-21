import type { ReadyState } from "react-use-websocket";

export interface WebSocketStatus {
  status: ReadyState;
  isAuthenticated: boolean;
}
