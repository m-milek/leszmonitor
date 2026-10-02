import { MonitorsApi } from "@/features/monitors/monitors-api";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { MonitorDetailPage } from "@/features/monitors/pages/MonitorDetailPage";
import { QUERY_KEYS } from "@/lib/consts";
import {
  DEFAULT_RESULTS_RANGE,
  parseResultsRange,
} from "@/features/monitors/components/DetailsPage/ResultsRangeSelect";

export const Route = createFileRoute("/_authenticated/monitors/$monitorSlug/")({
  validateSearch: (search): { range?: number } => ({
    range: parseResultsRange(search.range),
  }),
  beforeLoad: ({ params, search }) => {
    if (search.range === undefined) {
      throw redirect({
        to: "/monitors/$monitorSlug",
        params,
        search: { range: DEFAULT_RESULTS_RANGE },
        replace: true,
      });
    }
  },
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData({
      queryKey: [QUERY_KEYS.MONITORS, params.monitorSlug],
      queryFn: () => MonitorsApi.getBySlug(params.monitorSlug),
    }),
  head: ({ loaderData }) => ({
    meta: [{ title: `${loaderData?.name ?? "Monitor"} | Leszmonitor` }],
  }),
  component: RouteComponent,
});

function RouteComponent() {
  const { monitorSlug } = Route.useParams();
  const { range } = Route.useSearch();
  const navigate = Route.useNavigate();

  return (
    <MonitorDetailPage
      monitorSlug={monitorSlug}
      range={range ?? DEFAULT_RESULTS_RANGE}
      onRangeChange={(range) => navigate({ search: { range }, replace: true })}
    />
  );
}
