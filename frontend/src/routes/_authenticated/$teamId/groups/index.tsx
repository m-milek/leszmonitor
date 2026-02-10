import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { useQuery } from "@tanstack/react-query";
import { fetchGroups } from "@/lib/fetchGroups.ts";

export const Route = createFileRoute("/_authenticated/$teamId/groups/")({
  component: Groups,
});

function Groups() {
  const teamId = Route.useParams().teamId;

  const { data } = useQuery({
    queryKey: ["groups", teamId],
    queryFn: fetchGroups.bind(null, teamId),
  });

  if (!data) {
    return null;
  }

  return (
    <MainPanelContainer>
      <h1 className="text-2xl font-bold">Groups</h1>
      {data.length === 0 ? (
        <p className="mt-4 text-gray-500">No groups found.</p>
      ) : (
        <ul className="mt-4 space-y-2">
          {data.map((group) => (
            <li key={group.id} className="rounded-md bg-white p-4 shadow-sm">
              <h2 className="text-lg font-semibold">{group.name}</h2>
              <p className="text-sm text-gray-600">{group.description}</p>
            </li>
          ))}
        </ul>
      )}
    </MainPanelContainer>
  );
}
