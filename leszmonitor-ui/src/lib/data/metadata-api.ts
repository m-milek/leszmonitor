import { SERVER_API_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";

export interface MetadataResponse {
  ciBuildNumber: string;
  gitCommit: string;
  imageTag: string;
  version: string;
}

export const getMetadata = async (): Promise<MetadataResponse> => {
  const res = await authFetch(`${SERVER_API_URL}/instance-metadata`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return await res.json();
};
