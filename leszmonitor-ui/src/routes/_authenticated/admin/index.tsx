import { createFileRoute } from "@tanstack/react-router";
import { AdminPage } from "@/features/instance/pages/AdminPage";

export const Route = createFileRoute("/_authenticated/admin/")({
  head: () => ({
    meta: [{ title: "Administration | Leszmonitor" }],
  }),
  component: AdminPage,
});
