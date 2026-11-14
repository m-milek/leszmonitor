import { createFileRoute } from "@tanstack/react-router";
import { MonitorsPage } from "@/features/monitors/pages/MonitorsPage";
import { parseMonitorsSearch } from "@/features/monitors/monitors-search";

export const Route = createFileRoute("/_authenticated/monitors/")({
  validateSearch: parseMonitorsSearch,
  head: () => ({
    meta: [{ title: "Monitors | Leszmonitor" }],
  }),
  component: RouteComponent,
});

function RouteComponent() {
  const search = Route.useSearch();
  const navigate = Route.useNavigate();

  return (
    <MonitorsPage
      search={search}
      onSearchChange={(search) => navigate({ search, replace: true })}
    />
  );
}
