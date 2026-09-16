import { createFileRoute } from "@tanstack/react-router";
import { TagsPage } from "@/features/tags/pages/TagsPage.tsx";

export const Route = createFileRoute("/_authenticated/tags/")({
  head: () => ({
    meta: [{ title: "Tags | Leszmonitor" }],
  }),
  component: TagsPage,
});
