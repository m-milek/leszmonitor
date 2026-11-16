import { createFileRoute } from "@tanstack/react-router";
import { PublicStatusPage } from "@/features/status-pages/pages/PublicStatusPage";
import { StatusPagesApi } from "@/features/status-pages/status-pages-api";
import { QUERY_KEYS } from "@/lib/consts";

export const Route = createFileRoute("/status/$slug")({
  loader: ({ context, params }) =>
    context.queryClient
      .ensureQueryData({
        queryKey: [QUERY_KEYS.PUBLIC_STATUS_PAGES, params.slug],
        queryFn: () => StatusPagesApi.public.get(params.slug),
      })
      .catch(() => undefined),
  head: ({ loaderData }) => ({
    meta: [{ title: `${loaderData?.title ?? "Status"} | Leszmonitor` }],
  }),
  component: RouteComponent,
});

function RouteComponent() {
  const { slug } = Route.useParams();
  return <PublicStatusPage slug={slug} />;
}
