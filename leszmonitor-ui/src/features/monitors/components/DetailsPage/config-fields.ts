import type { ReactNode } from "react";

export type FieldConfig<T> = {
  [K in keyof T]?: {
    label: string;
    render?: (value: NonNullable<T[K]>, data: T) => ReactNode;
  };
};

export const formatValue = (value: unknown) =>
  typeof value === "string" ? value : JSON.stringify(value);

export const formatBoolean = (value: boolean) => (value ? "Yes" : "No");

export const formatList = (values: (string | number)[]) => values.join(", ");

export const formatMs = (value: number) => `${value} ms`;
