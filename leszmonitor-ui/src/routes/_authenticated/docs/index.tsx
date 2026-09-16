import { createFileRoute } from "@tanstack/react-router";
import { DocsPage } from "@/features/instance/pages/DocsPage";

export const Route = createFileRoute("/_authenticated/docs/")({
  head: () => ({
    meta: [{ title: "Documentation | Leszmonitor" }],
  }),
  component: DocsPage,
});
