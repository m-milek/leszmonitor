import { atom, createStore } from "jotai";
import type { Project, User } from "@/lib/types.ts";

export const store = createStore();

export const usernameAtom = atom<string | null>(null);
usernameAtom.debugLabel = "usernameAtom";

export const userAtom = atom<User | null>(null);
userAtom.debugLabel = "userAtom";

export const projectAtom = atom<Project | null>(null);
projectAtom.debugLabel = "projectAtom";
