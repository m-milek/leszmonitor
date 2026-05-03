export const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;
export const BACKEND_API_URL = `${BACKEND_URL}/api`;
export const BACKEND_WS_URL = BACKEND_API_URL.replace(/^http/, "ws");

export const QUERY_KEYS = {
  ORGS: "orgs",
  PROJECTS: "projects",
  USERS: "users",
  MONITORS: "monitors",
  MONITOR_RESULTS: "monitorResults",
};
