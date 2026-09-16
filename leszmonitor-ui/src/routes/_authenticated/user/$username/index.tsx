import { createFileRoute } from "@tanstack/react-router";
import { UserProfilePage } from "@/features/users/pages/UserProfilePage";

export const Route = createFileRoute("/_authenticated/user/$username/")({
  head: ({ params }) => ({
    meta: [{ title: `${params.username} | Leszmonitor` }],
  }),
  component: RouteComponent,
});

function RouteComponent() {
  const { username } = Route.useParams();
  return <UserProfilePage username={username} />;
}
