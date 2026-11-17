import { SERVER_API_URL } from "@/lib/consts";
import { authFetch } from "@/lib/api-client";

const getSpec = async (): Promise<Record<string, unknown>> => {
  const res = await authFetch(`${SERVER_API_URL}/swagger.json`, {
    method: "GET",
  });

  return await res.json();
};

export const DocsApi = {
  getSpec,
};
