import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/sidebar/Typography.tsx";
import { useQuery } from "@tanstack/react-query";
import { getGroup } from "@/lib/data/groupData.ts";

export const Route = createFileRoute(
  "/_authenticated/team/$teamId/groups/$groupId/",
)({
  component: GroupDetailsComponent,
});

function GroupDetailsComponent() {
  const { teamId, groupId } = Route.useParams();

  const { data } = useQuery({
    queryKey: ["groups", teamId, groupId],
    queryFn: () => getGroup(teamId, groupId),
  });

  if (!data) {
    return null;
  }

  return (
    <MainPanelContainer>
      <TypographyH1>Group Details</TypographyH1>
      <div>
        <p>
          <strong>Name:</strong> {data.name}
        </p>
        <p>
          <strong>Display ID:</strong> {data.displayId}
        </p>
        <p>
          <strong>Description:</strong> {data.description}
        </p>
        <p>
          <strong>Created At:</strong> {data.createdAt.toISOString()}
        </p>
        <p>
          <strong>Updated At:</strong> {data.updatedAt.toISOString()}
        </p>
      </div>
    </MainPanelContainer>
  );
}
