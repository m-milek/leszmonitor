import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export async function getLoginToken(): Promise<string | null> {
  const token = await cookieStore.get("LOGIN_TOKEN");
  return token?.value || null;
}
