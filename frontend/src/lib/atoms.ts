import { atom, createStore } from "jotai";
import type { Org, User } from "@/lib/types.ts";

export const store = createStore();

export const usernameAtom = atom<string | null>(null);
usernameAtom.debugLabel = "usernameAtom";

export const userAtom = atom<User | null>(null);
userAtom.debugLabel = "userAtom";

export const orgAtom = atom<Org | null>(null);
orgAtom.debugLabel = "orgAtom";
