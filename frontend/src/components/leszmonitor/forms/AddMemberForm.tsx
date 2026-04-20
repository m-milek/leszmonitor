import z from "zod";
import { useForm } from "@tanstack/react-form";
import { mapProjectRoleToDisplayName, ProjectRole } from "@/lib/types.ts";
import { Field, FieldLabel } from "@/components/ui/field.tsx";
import type { AddProjectMemberPayload } from "@/lib/data/projectData.ts";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";
import { LMCombobox } from "@/components/leszmonitor/forms/inputs/LMCombobox.tsx";
import { LMSelect } from "@/components/leszmonitor/forms/inputs/LMSelect.tsx";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/leszmonitor/forms/inputs/utils.ts";

const addMemberFormSchema = z.object({
  username: z.string().min(1, "Username is required"),
  role: z.enum(ProjectRole),
});

export interface AddMemberFormProps {
  onSubmitMember: (value: AddProjectMemberPayload) => Promise<unknown>;
  validUsernames: string[];
  formId?: string;
}

export function AddMemberForm({
  onSubmitMember,
  validUsernames,
  formId = "add-member-form",
}: AddMemberFormProps) {
  const roles = Object.values(ProjectRole);

  const form = useForm({
    defaultValues: {
      username: "",
      role: ProjectRole.Member,
    } as AddProjectMemberPayload,
    validators: {
      onSubmit: addMemberFormSchema,
    },
    onSubmit: async ({ value }) => {
      await onSubmitMember(value);
      form.reset();
    },
  });

  const roleSelectItems = roles.map((role) => ({
    value: role,
    label: mapProjectRoleToDisplayName[role],
  }));

  return (
    <form
      id={formId}
      onSubmit={(e) => {
        e.preventDefault();
        form.handleSubmit();
      }}
    >
      <Flex direction="column" className="gap-2">
        <form.Field
          name="username"
          children={(field) => {
            return (
              <Field>
                <FieldLabel htmlFor={field.name}>Username</FieldLabel>
                <LMCombobox
                  items={validUsernames}
                  value={field.state.value}
                  onValueChange={(value) => field.handleChange(value ?? "")}
                  placeholder="Find by username..."
                  id={field.name}
                  name={field.name}
                  className="max-w-1/2"
                  isInvalid={isFieldInvalid(field)}
                  errorMessage={getFirstError(field)}
                />
              </Field>
            );
          }}
        />
        <form.Field
          name="role"
          children={(field) => {
            return (
              <Field>
                <FieldLabel htmlFor={field.name}>Role</FieldLabel>
                <LMSelect
                  value={field.state.value}
                  onValueChange={(value) =>
                    field.handleChange(value as ProjectRole)
                  }
                  placeholder="Choose a role..."
                  items={roleSelectItems}
                  id={field.name}
                  name={field.name}
                  className="max-w-1/2"
                  isInvalid={isFieldInvalid(field)}
                  errorMessage={getFirstError(field)}
                />
              </Field>
            );
          }}
        />
      </Flex>
    </form>
  );
}
