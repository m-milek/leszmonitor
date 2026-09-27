import { SERVER_API_URL } from "@/lib/consts";
import { authFetch } from "@/lib/api-client";

export type GlobalParameterType = "string" | "int" | "bool" | "enum";

export type GlobalParameterSource = "config" | "ui" | "default";

export type GlobalParameterValue = string | number | boolean;

export interface GlobalParameter {
  key: string;
  displayName: string;
  description: string;
  type: GlobalParameterType;
  defaultValue: GlobalParameterValue;
  options?: string[];
  value: GlobalParameterValue;
  source: GlobalParameterSource;
}

const getAll = async (): Promise<GlobalParameter[]> => {
  const res = await authFetch(`${SERVER_API_URL}/global-parameters`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return (await res.json()) as GlobalParameter[];
};

const update = async (
  values: Record<string, GlobalParameterValue | null>,
): Promise<GlobalParameter[]> => {
  const res = await authFetch(`${SERVER_API_URL}/global-parameters`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(values),
  });

  return (await res.json()) as GlobalParameter[];
};

export const GlobalParametersApi = {
  getAll,
  update,
};
