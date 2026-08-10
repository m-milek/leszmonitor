import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useAppStore } from "@/lib/store.ts";
import { useEffect } from "react";
import { getProjectBySlug } from "@/lib/data/projectData.ts";
import { useQuery } from "@tanstack/react-query";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { PageContainer } from "@/components/leszmonitor/PageContainer.tsx";

export const Route = createFileRoute("/_authenticated/projects/$projectId")({
  component: ProjectLayout,
  loader: async ({ params, context }) => {
    const { projectId } = params;
    if (!projectId) {
      throw new Response("Project ID is required", { status: 400 });
    }

    try {
      return await context.queryClient.ensureQueryData({
        queryKey: ["project", projectId],
        queryFn: () => getProjectBySlug(projectId),
      });
    } catch {
      throw new Response("Project not found", { status: 404 });
    }
  },
  notFoundComponent: NotFound,
});

function ProjectLayout() {
  const { projectId } = Route.useParams();
  const { setProject } = useAppStore();

  const { data: project } = useQuery({
    queryKey: ["project", projectId],
    queryFn: () => getProjectBySlug(projectId),
  });

  useEffect(() => {
    if (project) {
      setProject(project);
    }
  }, [project, projectId, setProject]);

  return <Outlet />;
}

function NotFound() {
  return (
    <PageContainer>
      <TypographyH1>Not Found</TypographyH1>
    </PageContainer>
  );
}
