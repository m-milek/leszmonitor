import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import {
  formatDuration as dateFnsFormatDuration,
  intervalToDuration,
} from "date-fns";

export const cn = (...inputs: ClassValue[]): string => {
  return twMerge(clsx(inputs));
};

export const formatDate = (date: Date): string => {
  return date.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const DURATION_UNITS = [
  "years",
  "months",
  "days",
  "hours",
  "minutes",
  "seconds",
] as const;

export const formatDuration = (seconds: number): string => {
  const duration = intervalToDuration({ start: 0, end: seconds * 1000 });
  const format = DURATION_UNITS.filter((unit) => duration[unit]).slice(0, 2);
  return dateFnsFormatDuration(duration, { format });
};
