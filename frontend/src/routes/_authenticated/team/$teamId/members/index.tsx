import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { useAtomValue } from "jotai";
import { teamAtom } from "@/lib/atoms.ts";
import {
  TypographyH1,
  TypographyH2,
} from "@/components/leszmonitor/sidebar/Typography.tsx";
import { TeamMembersTable } from "@/components/leszmonitor/tables/TeamMembersTable.tsx";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Input } from "@/components/ui/input.tsx";
import { useState } from "react";
import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "@/components/ui/combobox.tsx";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx";
import { TeamRole } from "@/lib/types.ts";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  addUserToTeam,
  type AddUserToTeamPayload,
  fetchAllUsers,
  removeUserFromTeam,
} from "@/lib/data/userData.ts";
import { useForm } from "@tanstack/react-form";
import { z } from "zod";
import { Field, FieldError, FieldLabel } from "@/components/ui/field.tsx";

export const Route = createFileRoute("/_authenticated/team/$teamId/members/")({
  component: TeamRoute,
});

function TeamRoute() {
  const team = useAtomValue(teamAtom);

  const roles = Object.values(TeamRole);

  const queryClient = useQueryClient();

  const { data: users } = useQuery({
    queryKey: ["users"],
    queryFn: () => fetchAllUsers(),
  });

  const addMemberMutation = useMutation({
    mutationFn: (value: AddUserToTeamPayload) =>
      addUserToTeam(team!.displayId, {
        username: value.username,
        role: value.role,
      }),
    onSuccess: () => {
      form.reset();
      queryClient.invalidateQueries({ queryKey: ["team", team!.displayId] });
    },
  });

  const removeMemberMutation = useMutation({
    mutationFn: (username: string) =>
      removeUserFromTeam(team!.displayId, { username }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["team", team!.displayId] });
    },
  });
  const onMemberRemoved = async (username: string) => {
    console.log("Removing member", username);
    await removeMemberMutation.mutateAsync(username);
  };

  const addUserToTeamFormSchema = z.object({
    username: z.string().min(1, "Username is required"),
    role: z.enum(roles),
  });

  const form = useForm({
    defaultValues: {
      username: "",
      role: TeamRole.Member,
    } as AddUserToTeamPayload,
    validators: {
      onSubmit: addUserToTeamFormSchema,
    },
    onSubmit: async ({ value }) => {
      await addMemberMutation.mutateAsync(value);
    },
  });

  const [searchTerm, setSearchTerm] = useState("");

  if (!team || !users) {
    return null;
  }

  const validUsernames = users
    .map((user) => user.username)
    .filter((username) => {
      return !team.members.some((member) => member.username === username);
    });

  return (
    <MainPanelContainer>
      <TypographyH1>Manage Members</TypographyH1>
      <Card>
        <CardHeader>
          <TypographyH2>Add Members</TypographyH2>
        </CardHeader>
        <CardContent className="flex flex-col gap-6">
          <form
            id="add-member-form"
            onSubmit={(e) => {
              e.preventDefault();
              form.handleSubmit();
            }}
            className="flex items-end gap-4"
          >
            <form.Field
              name="username"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;
                return (
                  <Field>
                    <FieldLabel htmlFor={field.name}>Username</FieldLabel>
                    <Combobox
                      items={validUsernames}
                      value={field.state.value}
                      onValueChange={(value) => field.handleChange(value ?? "")}
                    >
                      <ComboboxInput
                        placeholder="Find by username..."
                        id={field.name}
                        name={field.name}
                      />
                      <ComboboxContent>
                        <ComboboxEmpty>No users found.</ComboboxEmpty>
                        <ComboboxList>
                          {(value) => {
                            return (
                              <ComboboxItem key={value} value={value}>
                                {value}
                              </ComboboxItem>
                            );
                          }}
                        </ComboboxList>
                      </ComboboxContent>
                    </Combobox>
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                );
              }}
            />
            <form.Field
              name="role"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;
                return (
                  <Field>
                    <FieldLabel htmlFor={field.name}>Role</FieldLabel>
                    <Select
                      onValueChange={(value) =>
                        field.handleChange(value as TeamRole)
                      }
                      defaultValue={field.state.value}
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="Choose a role..." />
                      </SelectTrigger>
                      <SelectContent>
                        {roles.map((role) => (
                          <SelectItem key={role} value={role}>
                            {role}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                );
              }}
            />
            <Button type="submit">Add Member</Button>
          </form>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <div className="flex justify-between">
            <TypographyH2>
              {team.members.length}{" "}
              {team.members.length === 1 ? "Member" : "Members"}
            </TypographyH2>
          </div>
          <Input
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder="Search members..."
            className="w-[50%]"
          />
        </CardHeader>
        <CardContent>
          <TeamMembersTable
            onMemberRemoved={onMemberRemoved}
            members={team.members.filter((member) =>
              member.username.toLowerCase().includes(searchTerm.toLowerCase()),
            )}
          />
        </CardContent>
        <CardFooter></CardFooter>
      </Card>
    </MainPanelContainer>
  );
}
