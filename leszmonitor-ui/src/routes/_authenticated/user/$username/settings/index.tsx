import { createFileRoute } from "@tanstack/react-router";
import { UserSettingsPage } from "@/features/users/pages/UserSettingsPage.tsx";

export const Route = createFileRoute(
  "/_authenticated/user/$username/settings/",
)({
  head: () => ({
    meta: [{ title: "Settings | Leszmonitor" }],
  }),
  component: UserSettingsPage,
});
