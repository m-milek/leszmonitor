import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  TypographyH1,
  TypographyH2,
} from "@/components/leszmonitor/Typography.tsx";
import {
  addGroup,
  deleteGroup,
  getGroups,
  type GroupInput,
} from "@/lib/data/groupData.ts";
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
import { GroupsTable } from "@/components/leszmonitor/tables/GroupsTable.tsx";

export const Route = createFileRoute("/_authenticated/team/$teamId/groups/")({
  component: Groups,
});

const groupFormSchema = z.object({
  name: z.string().min(1, "Group name is required"),
  displayId: z.string().min(1, "Display ID is required"),
  description: z.string(),
});

function Groups() {
  const teamId = Route.useParams().teamId;
  const queryClient = useQueryClient();

  const { data } = useQuery({
    queryKey: ["groups", teamId],
    queryFn: () => getGroups(teamId),
  });

  const addGroupMutation = useMutation({
    mutationFn: (newGroup: GroupInput) => addGroup(teamId, newGroup),
    onSuccess: () => {
      form.reset();
      queryClient.invalidateQueries({ queryKey: ["groups", teamId] });
    },
  });

  const deleteGroupMutation = useMutation({
    mutationFn: (groupId: string) => deleteGroup(teamId, groupId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["groups", teamId] });
    },
  });

  const form = useForm({
    defaultValues: {
      name: "",
      displayId: "",
      description: "",
    },
    validators: {
      onSubmit: groupFormSchema,
    },
    onSubmit: async ({ value }) => {
      await addGroupMutation.mutateAsync(value);
    },
  });

  if (!data) {
    return null;
  }

  return (
    <MainPanelContainer>
      <TypographyH1>Groups</TypographyH1>
      <Card>
        <CardHeader>
          <TypographyH2>Create New Group</TypographyH2>
        </CardHeader>
        <CardContent>
          <form
            id="group-form"
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
                        <FieldLabel htmlFor={field.name}>Group Name</FieldLabel>
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
            form="group-form"
            disabled={addGroupMutation.isPending}
          >
            {addGroupMutation.isPending ? "Adding..." : "Add Group"}
          </Button>
        </CardFooter>
      </Card>
      <Card>
        <CardHeader>
          <TypographyH2>Existing Groups</TypographyH2>
        </CardHeader>
        <CardContent>
          <GroupsTable
            groups={data}
            teamId={teamId}
            onGroupDeleted={async (groupId) =>
              deleteGroupMutation.mutateAsync(groupId)
            }
          />
        </CardContent>
      </Card>
    </MainPanelContainer>
  );
}
