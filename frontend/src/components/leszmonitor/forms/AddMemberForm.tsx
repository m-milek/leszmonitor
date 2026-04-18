import z from "zod";
import { useForm } from "@tanstack/react-form";
import { mapOrgRoleToDisplayName, OrgRole } from "@/lib/types.ts";
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
import { Field, FieldError, FieldLabel } from "@/components/ui/field.tsx";
import type { AddUserToOrgPayload } from "@/lib/data/userData.ts";
import { Flex } from "@/components/leszmonitor/ui/Flex.tsx";

const addUserToOrgFormSchema = z.object({
  username: z.string().min(1, "Username is required"),
  role: z.nativeEnum(OrgRole),
});

export interface AddMemberFormProps {
  onSubmitMember: (value: AddUserToOrgPayload) => Promise<unknown>;
  validUsernames: string[];
  formId?: string;
}

export function AddMemberForm({
  onSubmitMember,
  validUsernames,
  formId = "add-member-form",
}: AddMemberFormProps) {
  const roles = Object.values(OrgRole);

  const form = useForm({
    defaultValues: {
      username: "",
      role: OrgRole.Member,
    } as AddUserToOrgPayload,
    validators: {
      onSubmit: addUserToOrgFormSchema,
    },
    onSubmit: async ({ value }) => {
      await onSubmitMember(value);
      form.reset();
    },
  });

  return (
    <form
      id={formId}
      onSubmit={(e) => {
        e.preventDefault();
        form.handleSubmit();
      }}
    >
      <Flex direction="vertical" gap="0.5rem">
        <form.Field
          name="username"
          children={(field) => {
            const isInvalid =
              field.state.meta.isTouched && !field.state.meta.isValid;
            return (
              <Field>
                <FieldLabel htmlFor={field.name}>Username</FieldLabel>
                <Flex direction="horizontal">
                  <Combobox
                    items={validUsernames}
                    value={field.state.value}
                    onValueChange={(value) => field.handleChange(value ?? "")}
                  >
                    <ComboboxInput
                      placeholder="Find by username..."
                      id={field.name}
                      name={field.name}
                      className="max-w-1/2"
                      autoComplete="off"
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
                  {isInvalid && <FieldError errors={field.state.meta.errors} />}
                </Flex>
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
                <Flex direction="horizontal">
                  <Select
                    onValueChange={(value) =>
                      field.handleChange(value as OrgRole)
                    }
                    defaultValue={field.state.value}
                  >
                    <SelectTrigger className="max-w-1/2">
                      <SelectValue placeholder="Choose a role..." />
                    </SelectTrigger>
                    <SelectContent>
                      {roles.map((role) => (
                        <SelectItem key={role} value={role}>
                          {mapOrgRoleToDisplayName[role]}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {isInvalid && <FieldError errors={field.state.meta.errors} />}
                </Flex>
              </Field>
            );
          }}
        />
      </Flex>
    </form>
  );
}
