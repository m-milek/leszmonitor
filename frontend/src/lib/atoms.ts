import { atom, createStore } from "jotai";
import type { Project, User } from "@/lib/types.ts";
import { ReadyState } from "react-use-websocket";

export const store = createStore();

export const usernameAtom = atom<string | null>(null);
usernameAtom.debugLabel = "usernameAtom";

export const userAtom = atom<User | null>(null);
userAtom.debugLabel = "userAtom";

export const projectAtom = atom<Project | null>(null);
projectAtom.debugLabel = "projectAtom";

export const webSocketConnectionStatusAtom = atom<ReadyState>(
  ReadyState.CLOSED,
);
webSocketConnectionStatusAtom.debugLabel = "webSocketAtom";
