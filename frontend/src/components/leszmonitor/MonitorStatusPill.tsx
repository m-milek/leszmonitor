import type { Monitor, MonitorState } from "@/lib/types.ts";

export interface MonitorStatusPillProps {
  monitor: Monitor;
}

const mapMonitorState = (state: MonitorState) => {
  switch (state) {
    case "active":
      return {
        text: "Active",
        color: "bg-green-500 text-white",
      };
    case "paused":
      return {
        text: "Paused",
        color: "bg-gray-500 text-white",
      };
    default:
      return {
        text: "Invalid",
        color: "bg-gray-300 text-gray-700",
      };
  }
};

export const MonitorStatusPill = ({ monitor }: MonitorStatusPillProps) => {
  const text = mapMonitorState(monitor.state).text;
  const color = mapMonitorState(monitor.state).color;

  return (
    <span
      className={`inline-flex items-center px-4 py-1 rounded-xl ${color} font-medium`}
    >
      {text}
    </span>
  );
};
