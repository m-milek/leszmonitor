import { z } from "zod";

export const stringifySearch = (search: Record<string, unknown>): string => {
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(search)) {
    for (const item of [value].flat()) {
      if (item !== undefined && item !== null) params.append(key, String(item));
    }
  }
  const query = params.toString();
  return query ? `?${query}` : "";
};

export const parseSearch = (searchStr: string): Record<string, unknown> => {
  const params = new URLSearchParams(searchStr);
  return Object.fromEntries(
    [...new Set(params.keys())].map((key) => {
      const values = params.getAll(key);
      return [key, values.length > 1 ? values : values[0]];
    }),
  );
};

export const searchArray = <T extends z.ZodType>(item: T) =>
  z
    .preprocess(
      (value) => (value === undefined ? undefined : [value].flat()),
      z.array(item).optional(),
    )
    .catch(undefined);
