import { createFileRoute } from "@tanstack/react-router";
import { StatusPagesPage } from "@/features/status-pages/pages/StatusPagesPage";

export const Route = createFileRoute("/_authenticated/statusPages/")({
  head: () => ({
    meta: [{ title: "Status Pages | Leszmonitor" }],
  }),
  component: StatusPagesPage,
});
