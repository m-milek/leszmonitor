const SERVER_HOSTNAME = "localhost:7001";
export const SERVER_API_URL = `http://${SERVER_HOSTNAME}/api/v1`;
export const SERVER_WS_URL = `ws://${SERVER_HOSTNAME}/api/ws`;

export const QUERY_KEYS = {
  ORGS: "orgs",
  PROJECTS: "projects",
  USERS: "users",
  MONITORS: "monitors",
  MONITOR_RESULTS: "monitorResults",
  MONITOR_AVERAGE_LATENCY: "monitorAverageLatency"
};
