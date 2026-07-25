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
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { NewProjectForm } from "@/components/leszmonitor/forms/NewProjectForm.tsx";
import { ProjectsTable } from "@/components/leszmonitor/tables/ProjectsTable.tsx";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx";

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
          <Flex direction="row" className="justify-between items-center">
            <TypographyH2>Your Projects</TypographyH2>
            <Dialog>
              <DialogTrigger asChild>
                <Button>Create Project</Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create New Project</DialogTitle>
                </DialogHeader>
                <NewProjectForm
                  formId="project-form"
                  onSubmitProject={async (value) => {
                    await addProjectMutation.mutateAsync(value);
                  }}
                />
                <DialogFooter>
                  <Button
                    type="submit"
                    form="project-form"
                    disabled={addProjectMutation.isPending}
                  >
                    {addProjectMutation.isPending
                      ? "Adding..."
                      : "Create Project"}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          </Flex>
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
