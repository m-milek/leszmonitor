import { SERVER_API_URL } from "@/lib/consts";
import { authFetch } from "@/lib/api-client";

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
