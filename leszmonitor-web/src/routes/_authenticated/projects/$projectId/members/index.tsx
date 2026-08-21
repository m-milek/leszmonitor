import { createFileRoute } from "@tanstack/react-router";
import { PageContainer } from "@/components/leszmonitor/PageContainer.tsx";
import { useAppStore } from "@/lib/store.ts";
import {
  TypographyH1,
  TypographyH2,
} from "@/components/leszmonitor/ui/Typography.tsx";
import { ProjectMembersTable } from "@/components/leszmonitor/tables/OrgMembersTable.tsx";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Input } from "@/components/ui/input.tsx";
import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  addMemberToProject,
  type AddProjectMemberPayload,
  removeMemberFromProject,
} from "@/lib/data/projectData.ts";
import { getAllUsers } from "@/lib/data/userData.ts";
import { AddMemberForm } from "@/components/leszmonitor/forms/AddMemberForm.tsx";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx";

export const Route = createFileRoute(
  "/_authenticated/projects/$projectId/members/",
)({
  component: ProjectMembersRoute,
});

function ProjectMembersRoute() {
  const { project } = useAppStore();
  const queryClient = useQueryClient();

  const { data: users } = useQuery({
    queryKey: ["users"],
    queryFn: () => getAllUsers(),
  });

  const addMemberMutation = useMutation({
    mutationFn: (value: AddProjectMemberPayload) =>
      addMemberToProject(project!.slug, value),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["project", project!.slug] });
    },
  });

  const removeMemberMutation = useMutation({
    mutationFn: (username: string) =>
      removeMemberFromProject(project!.slug, { username }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["project", project!.slug] });
    },
  });

  const onMemberRemoved = async (username: string) => {
    await removeMemberMutation.mutateAsync(username);
  };

  const [searchTerm, setSearchTerm] = useState("");

  if (!project || !users) {
    return null;
  }

  const validUsernames = users
    .map((user) => user.username)
    .filter((username) => {
      return !project.members.some((member) => member.username === username);
    });

  return (
    <PageContainer>
      <TypographyH1>Manage Members</TypographyH1>
      <Card>
        <CardHeader>
          <Flex direction="row" className="justify-between">
            <TypographyH2>
              {project.members.length}{" "}
              {project.members.length === 1 ? "Member" : "Members"}
            </TypographyH2>
            <Dialog>
              <DialogTrigger asChild>
                <Button>Add Member</Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Add Member</DialogTitle>
                </DialogHeader>
                <AddMemberForm
                  formId="add-member-form"
                  onSubmitMember={async (value) => {
                    await addMemberMutation.mutateAsync(value);
                  }}
                  validUsernames={validUsernames}
                />
                <DialogFooter>
                  <Button
                    type="submit"
                    form="add-member-form"
                    disabled={addMemberMutation.isPending}
                  >
                    {addMemberMutation.isPending ? "Adding..." : "Submit"}
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          </Flex>
          <Input
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder="Search members..."
            className="w-[50%]"
          />
        </CardHeader>
        <CardContent>
          <ProjectMembersTable
            onMemberRemoved={onMemberRemoved}
            members={project.members.filter((member) =>
              member.username.toLowerCase().includes(searchTerm.toLowerCase()),
            )}
          />
        </CardContent>
        <CardFooter></CardFooter>
      </Card>
    </PageContainer>
  );
}
