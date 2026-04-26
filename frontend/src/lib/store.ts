import { create } from 'zustand';
import type { Project, User } from '@/lib/types.ts';
import { ReadyState } from 'react-use-websocket';
import type { WebSocketStatus } from '@/lib/data/webSocket.ts';

interface AppState {
  username: string | null;
  setUsername: (username: string | null) => void;
  user: User | null;
  setUser: (user: User | null) => void;
  project: Project | null;
  setProject: (project: Project | null) => void;
  webSocketConnectionStatus: WebSocketStatus;
  setWebSocketConnectionStatus: (status: WebSocketStatus) => void;
}

export const useAppStore = create<AppState>((set) => ({
  username: null,
  setUsername: (username) => set({ username }),
  user: null,
  setUser: (user) => set({ user }),
  project: null,
  setProject: (project) => set({ project }),
  webSocketConnectionStatus: {
    status: ReadyState.CLOSED,
    isAuthenticated: false,
  } as WebSocketStatus,
  setWebSocketConnectionStatus: (webSocketConnectionStatus) =>
    set({ webSocketConnectionStatus }),
}));
