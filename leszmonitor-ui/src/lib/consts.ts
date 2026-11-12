const SERVER_HOSTNAME = import.meta.env.DEV
  ? "localhost:7001"
  : window.location.host;
const SECURE = !import.meta.env.DEV && window.location.protocol === "https:";

export const SERVER_API_URL = `${SECURE ? "https" : "http"}://${SERVER_HOSTNAME}/api/v1`;
export const SERVER_WS_URL = `${SECURE ? "wss" : "ws"}://${SERVER_HOSTNAME}/api/v1/ws`;

export const QUERY_KEYS = {
  USERS: "users",
  MONITORS: "monitors",
  MONITOR_RESULTS: "monitorResults",
  MONITOR_LATENCY_STATS: "monitorLatencyStats",
  MONITOR_STATUS_HISTORY: "monitorStatusHistory",
  TAGS: "tags",
  GLOBAL_PARAMETERS: "globalParameters",
  API_SPEC: "apiSpec",
};
