export interface MetadataResponse {
  ciBuildNumber: string;
  gitCommit: string;
  imageTag: string;
  version: string;
}

export const getMetadata = async (): Promise<MetadataResponse> => {
  // TODO call real API
  return {
    ciBuildNumber: "#67",
    gitCommit: "abcd123123321",
    imageTag: "leszmonitor:1.0.0",
    version: "v1.0.0",
  };
};
