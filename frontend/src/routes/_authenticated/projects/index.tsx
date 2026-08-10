import { createFileRoute } from "@tanstack/react-router";
import { PageContainer } from "@/components/leszmonitor/PageContainer.tsx";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  TypographyH1,
  TypographyH2,
} from "@/components/leszmonitor/ui/Typography.tsx";
import {
  addProject,
  getProjects,
  type ProjectInput,
} from "@/lib/data/projectData.ts";
import { QUERY_KEYS } from "@/lib/consts.ts";
import { ProjectListItem } from "@/components/leszmonitor/ProjectListItem.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx";
import { NewProjectForm } from "@/components/leszmonitor/forms/NewProjectForm";
import { useState } from "react";

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

  const [formOpen, setFormOpen] = useState(false);

  if (!data) {
    return null;
  }

  return (
    <PageContainer>
      <TypographyH1>Projects</TypographyH1>
      <Card>
        <CardHeader>
          <Flex direction="row" className="justify-between items-center">
            <TypographyH2>Your Projects</TypographyH2>
            <Dialog open={formOpen} onOpenChange={setFormOpen}>
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
                    setFormOpen(false);
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
          <Flex direction="column" className="gap-4">
            {data.map((project) => (
              <div key={project.id}>
                <ProjectListItem project={project} />
              </div>
            ))}
          </Flex>
        </CardContent>
      </Card>
    </PageContainer>
  );
}
