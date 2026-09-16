import { createFileRoute } from "@tanstack/react-router";
import { NewMonitorPage } from "@/features/monitors/pages/NewMonitorPage";

export const Route = createFileRoute("/_authenticated/monitors/new/")({
  head: () => ({
    meta: [{ title: "New Monitor | Leszmonitor" }],
  }),
  component: NewMonitorPage,
});
