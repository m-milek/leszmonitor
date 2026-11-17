import { createFileRoute } from "@tanstack/react-router";
import { AdminPage } from "@/features/settings/pages/AdminPage";

export const Route = createFileRoute("/_authenticated/settings/")({
  head: () => ({
    meta: [{ title: "Administration | Leszmonitor" }],
  }),
  component: AdminPage,
});
