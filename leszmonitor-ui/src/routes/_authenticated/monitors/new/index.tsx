import { createFileRoute } from "@tanstack/react-router";
import { NewMonitorPage } from "@/features/monitors/pages/NewMonitorPage.tsx";

export const Route = createFileRoute("/_authenticated/monitors/new/")({
  head: () => ({
    meta: [{ title: "New Monitor | Leszmonitor" }],
  }),
  component: NewMonitorPage,
});
