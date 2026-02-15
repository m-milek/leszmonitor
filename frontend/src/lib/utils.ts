import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { TeamRole } from "@/lib/types.ts";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export async function getLoginToken(): Promise<string | null> {
  const token = await cookieStore.get("LOGIN_TOKEN");
  return token?.value || null;
}

export function formatDate(date: Date): string {
  return date.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function formatRole(role: TeamRole): string {
  switch (role) {
    case TeamRole.Viewer:
      return "Viewer";
    case TeamRole.Member:
      return "Member";
    case TeamRole.Admin:
      return "Admin";
    case TeamRole.Owner:
      return "Owner";
    default:
      return role;
  }
}
