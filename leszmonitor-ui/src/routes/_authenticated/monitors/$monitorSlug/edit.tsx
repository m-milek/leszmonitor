import { createFileRoute } from "@tanstack/react-router";
import { getMonitorBySlug } from "@/features/monitors/monitors-api.ts";
import { MonitorEditPage } from "@/features/monitors/pages/MonitorEditPage.tsx";
import { QUERY_KEYS } from "@/lib/consts.ts";

export const Route = createFileRoute(
  "/_authenticated/monitors/$monitorSlug/edit",
)({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData({
      queryKey: [QUERY_KEYS.MONITORS, params.monitorSlug],
      queryFn: () => getMonitorBySlug(params.monitorSlug),
    }),
  head: ({ loaderData }) => ({
    meta: [{ title: `Edit ${loaderData?.name ?? "Monitor"} | Leszmonitor` }],
  }),
  component: RouteComponent,
});

function RouteComponent() {
  const { monitorSlug } = Route.useParams();
  return <MonitorEditPage monitorSlug={monitorSlug} />;
}
