import { createFileRoute } from "@tanstack/react-router";
import { AdminPage } from "@/features/users/pages/AdminPage.tsx";

export const Route = createFileRoute("/_authenticated/admin/")({
  head: () => ({
    meta: [{ title: "Administration | Leszmonitor" }],
  }),
  component: AdminPage,
});
