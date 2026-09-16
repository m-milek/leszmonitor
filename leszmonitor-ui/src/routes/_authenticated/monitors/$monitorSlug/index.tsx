import { createFileRoute } from "@tanstack/react-router";
import { getMonitorBySlug } from "@/features/monitors/monitors-api";
import { MonitorDetailPage } from "@/features/monitors/pages/MonitorDetailPage";
import { QUERY_KEYS } from "@/lib/consts";

export const Route = createFileRoute("/_authenticated/monitors/$monitorSlug/")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData({
      queryKey: [QUERY_KEYS.MONITORS, params.monitorSlug],
      queryFn: () => getMonitorBySlug(params.monitorSlug),
    }),
  head: ({ loaderData }) => ({
    meta: [{ title: `${loaderData?.name ?? "Monitor"} | Leszmonitor` }],
  }),
  component: RouteComponent,
});

function RouteComponent() {
  const { monitorSlug } = Route.useParams();
  return <MonitorDetailPage monitorSlug={monitorSlug} />;
}
