import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { useQuery } from "@tanstack/react-query";
import { getProject } from "@/lib/data/projectData.ts";

export const Route = createFileRoute(
  "/_authenticated/org/$orgId/projects/$projectId/",
)({
  component: ProjectDetailsComponent,
});

function ProjectDetailsComponent() {
  const { orgId, projectId } = Route.useParams();

  const { data } = useQuery({
    queryKey: ["projects", orgId, projectId],
    queryFn: () => getProject(orgId, projectId),
  });

  if (!data) {
    return null;
  }

  return (
    <MainPanelContainer>
      <TypographyH1>Project Details</TypographyH1>
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
