import { create } from "zustand";
import type { User } from "@/features/users/types.ts";
import { ReadyState } from "react-use-websocket";
import type { WebSocketStatus } from "@/app/providers/websocket-status.ts";

interface AppState {
  username: string | null;
  setUsername: (username: string | null) => void;
  user: User | null;
  setUser: (user: User | null) => void;
  webSocketConnectionStatus: WebSocketStatus;
  setWebSocketConnectionStatus: (status: WebSocketStatus) => void;
}

export const useAppStore = create<AppState>((set) => ({
  username: null,
  setUsername: (username) => set({ username }),
  user: null,
  setUser: (user) => set({ user }),
  webSocketConnectionStatus: {
    status: ReadyState.CLOSED,
    isAuthenticated: false,
  } as WebSocketStatus,
  setWebSocketConnectionStatus: (webSocketConnectionStatus) =>
    set({ webSocketConnectionStatus }),
}));
