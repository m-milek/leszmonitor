import { createFileRoute } from "@tanstack/react-router";
import { MonitorsPage } from "@/features/monitors/pages/MonitorsPage";

export const Route = createFileRoute("/_authenticated/monitors/")({
  head: () => ({
    meta: [{ title: "Monitors | Leszmonitor" }],
  }),
  component: MonitorsPage,
});
