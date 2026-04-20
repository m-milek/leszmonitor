import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useSetAtom } from "jotai";
import { projectAtom } from "@/lib/atoms.ts";
import { useEffect } from "react";
import { getProject } from "@/lib/data/projectData.ts";
import { useQuery } from "@tanstack/react-query";
import { TypographyH1 } from "@/components/leszmonitor/ui/Typography.tsx";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";

export const Route = createFileRoute("/_authenticated/projects/$projectId")({
  component: ProjectLayout,
  notFoundComponent: NotFound,
});

function ProjectLayout() {
  const { projectId } = Route.useParams();
  const setProjectAtom = useSetAtom(projectAtom);

  const { data: project } = useQuery({
    queryKey: ["project", projectId],
    queryFn: () => getProject(projectId),
  });

  useEffect(() => {
    if (project) {
      setProjectAtom(project);
    }
  }, [project, projectId, setProjectAtom]);

  return <Outlet />;
}

function NotFound() {
  return (
    <MainPanelContainer>
      <TypographyH1>Not Found</TypographyH1>;
    </MainPanelContainer>
  );
}
