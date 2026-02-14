import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_authenticated/team/$teamId/")({
  component: TeamDashboard,
});

function TeamDashboard() {
  const { teamId } = Route.useParams();

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">Team Dashboard</h1>
      <p className="mt-2">
        Team: <span className="font-semibold">{teamId}</span>
      </p>
    </div>
  );
}
