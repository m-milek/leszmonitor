import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  TypographyH1,
  TypographyH2,
} from "@/components/leszmonitor/ui/Typography.tsx";
import {
  addProject,
  deleteProject,
  getProjects,
  type ProjectInput,
} from "@/lib/data/projectData.ts";
import { Button } from "@/components/ui/button.tsx";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card.tsx";
import { NewProjectForm } from "@/components/leszmonitor/forms/NewProjectForm.tsx";
import { ProjectsTable } from "@/components/leszmonitor/tables/ProjectsTable.tsx";
import { QUERY_KEYS } from "@/lib/consts.ts";

export const Route = createFileRoute("/_authenticated/projects/")({
  component: Projects,
});

function Projects() {
  const queryClient = useQueryClient();

  const { data } = useQuery({
    queryKey: [QUERY_KEYS.PROJECTS],
    queryFn: () => getProjects(),
  });

  const addProjectMutation = useMutation({
    mutationFn: (newProject: ProjectInput) => addProject(newProject),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.PROJECTS] });
    },
  });

  const deleteProjectMutation = useMutation({
    mutationFn: (projectId: string) => deleteProject(projectId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.PROJECTS] });
    },
  });

  if (!data) {
    return null;
  }

  return (
    <MainPanelContainer>
      <TypographyH1>Projects</TypographyH1>
      <Card>
        <CardHeader>
          <TypographyH2>Create New Project</TypographyH2>
        </CardHeader>
        <CardContent>
          <NewProjectForm
            formId="project-form"
            onSubmitProject={(value) => addProjectMutation.mutateAsync(value)}
          />
        </CardContent>
        <CardFooter className="justify-end">
          <Button
            type="submit"
            form="project-form"
            disabled={addProjectMutation.isPending}
          >
            {addProjectMutation.isPending ? "Adding..." : "Add Project"}
          </Button>
        </CardFooter>
      </Card>
      <Card>
        <CardHeader>
          <TypographyH2>Your Projects</TypographyH2>
        </CardHeader>
        <CardContent>
          <ProjectsTable
            projects={data}
            onProjectDeleted={async (projectId) =>
              deleteProjectMutation.mutateAsync(projectId)
            }
          />
        </CardContent>
      </Card>
    </MainPanelContainer>
  );
}
