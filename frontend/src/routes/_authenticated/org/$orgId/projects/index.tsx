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
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field.tsx";
import { Input } from "@/components/ui/input.tsx";
import { useForm } from "@tanstack/react-form";
import { z } from "zod";
import { Textarea } from "@/components/ui/textarea.tsx";
import { ProjectsTable } from "@/components/leszmonitor/tables/ProjectsTable.tsx";

export const Route = createFileRoute("/_authenticated/org/$orgId/projects/")({
  component: Projects,
});

const projectFormSchema = z.object({
  name: z.string().min(1, "Project name is required"),
  displayId: z.string().min(1, "Display ID is required"),
  description: z.string(),
});

function Projects() {
  const { orgId } = Route.useParams();
  const queryClient = useQueryClient();

  const { data } = useQuery({
    queryKey: ["projects", orgId],
    queryFn: () => getProjects(orgId),
  });

  const addProjectMutation = useMutation({
    mutationFn: (newProject: ProjectInput) => addProject(orgId, newProject),
    onSuccess: () => {
      form.reset();
      queryClient.invalidateQueries({ queryKey: ["projects", orgId] });
    },
  });

  const deleteProjectMutation = useMutation({
    mutationFn: (projectId: string) => deleteProject(orgId, projectId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", orgId] });
    },
  });

  const form = useForm({
    defaultValues: {
      name: "",
      displayId: "",
      description: "",
    },
    validators: {
      onSubmit: projectFormSchema,
    },
    onSubmit: async ({ value }) => {
      await addProjectMutation.mutateAsync(value);
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
          <form
            id="project-form"
            onSubmit={(e) => {
              e.preventDefault();
              form.handleSubmit();
            }}
          >
            <FieldGroup className="gap-2">
              <div className="flex gap-8">
                <form.Field
                  name="name"
                  children={(field) => {
                    const isInvalid =
                      field.state.meta.isTouched && !field.state.meta.isValid;
                    return (
                      <Field>
                        <FieldLabel htmlFor={field.name}>Project Name</FieldLabel>
                        <Input
                          id={field.name}
                          name={field.name}
                          value={field.state.value}
                          onChange={(e) => field.handleChange(e.target.value)}
                          autoComplete="off"
                        />
                        {isInvalid && (
                          <FieldError errors={field.state.meta.errors} />
                        )}
                      </Field>
                    );
                  }}
                />
                <form.Field
                  name="displayId"
                  children={(field) => {
                    const isInvalid =
                      field.state.meta.isTouched && !field.state.meta.isValid;
                    return (
                      <Field>
                        <FieldLabel htmlFor={field.name}>Display ID</FieldLabel>
                        <Input
                          id={field.name}
                          name={field.name}
                          value={field.state.value}
                          onChange={(e) => field.handleChange(e.target.value)}
                          autoComplete="off"
                        />
                        {isInvalid && (
                          <FieldError errors={field.state.meta.errors} />
                        )}
                      </Field>
                    );
                  }}
                />
              </div>
              <form.Field
                name="description"
                children={(field) => {
                  return (
                    <Field>
                      <FieldLabel htmlFor={field.name}>
                        Description (Optional)
                      </FieldLabel>
                      <Textarea
                        id={field.name}
                        name={field.name}
                        value={field.state.value}
                        onChange={(e) => field.handleChange(e.target.value)}
                        autoComplete="off"
                      />
                    </Field>
                  );
                }}
              />
            </FieldGroup>
          </form>
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
          <TypographyH2>Existing Projects</TypographyH2>
        </CardHeader>
        <CardContent>
          <ProjectsTable
            projects={data}
            orgId={orgId}
            onProjectDeleted={async (projectId) =>
              deleteProjectMutation.mutateAsync(projectId)
            }
          />
        </CardContent>
      </Card>
    </MainPanelContainer>
  );
}
