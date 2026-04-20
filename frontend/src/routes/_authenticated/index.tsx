import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_authenticated/")({
  beforeLoad: async () => {
    throw redirect({ to: "/projects" });
  },
  component: () => null,
});
