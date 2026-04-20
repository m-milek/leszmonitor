import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { ProjectRole } from "@/lib/types.ts";

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

export function formatRole(role: ProjectRole): string {
  switch (role) {
    case ProjectRole.Viewer:
      return "Viewer";
    case ProjectRole.Member:
      return "Member";
    case ProjectRole.Admin:
      return "Admin";
    case ProjectRole.Owner:
      return "Owner";
    default:
      return role;
  }
}
