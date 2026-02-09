import { atom } from "jotai";
import type { User } from "@/lib/types.ts";

export const usernameAtom = atom<string | null>(null);

export const userAtom = atom<User | null>(null);
userAtom.debugLabel = "userAtom";
