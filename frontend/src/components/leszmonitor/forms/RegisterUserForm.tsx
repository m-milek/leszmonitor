import { z } from "zod";
import { useForm } from "@tanstack/react-form";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field.tsx";
import { LMInputField } from "@/components/leszmonitor/forms/inputs/LMInputField.tsx";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/leszmonitor/forms/inputs/utils.ts";
import { type RegisterUserPayload } from "@/lib/data/userData.ts";

export interface RegisterUserFormProps {
  id?: string;
  requirePasswordConfirm?: boolean;
  onSubmit: (values: RegisterUserPayload) => Promise<void>;
  className?: string;
}

export function RegisterUserForm({
  id = "register-form",
  requirePasswordConfirm = true,
  onSubmit,
  className,
}: RegisterUserFormProps) {
  const formSchema = z
    .object({
      username: z
        .string()
        .min(2, "Username has to be at least 2 characters long"),
      password: z
        .string()
        .min(6, "Password has to be at least 6 characters long"),
      passwordConfirm: requirePasswordConfirm
        ? z.string().min(6, "Verify the password by entering it again")
        : z.string().optional(),
    })
    .refine(
      (data) => {
        if (!requirePasswordConfirm) return true;
        return data.password === data.passwordConfirm;
      },
      {
        message: "Passwords don't match",
        path: ["passwordConfirm"],
      }
    );

  const form = useForm({
    defaultValues: {
      username: "",
      password: "",
      passwordConfirm: "",
    },
    validators: {
      onSubmit: formSchema,
    },
    onSubmit: async ({ value }) => {
      await onSubmit({
        username: value.username,
        password: value.password,
      });
    },
  });

  return (
    <form
      id={id}
      className={className}
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}
    >
      <FieldGroup className="gap-2">
        <form.Field
          name="username"
          children={(field) => (
            <Field id={field.name}>
              <FieldLabel>Username</FieldLabel>
              <LMInputField
                name={field.name}
                type="text"
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                isInvalid={isFieldInvalid(field)}
                errorMessage={getFirstError(field)}
              />
            </Field>
          )}
        />
        <form.Field
          name="password"
          children={(field) => (
            <Field id={field.name}>
              <FieldLabel>Password</FieldLabel>
              <LMInputField
                name={field.name}
                type="password"
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                autoComplete="new-password"
                isInvalid={isFieldInvalid(field)}
                errorMessage={getFirstError(field)}
              />
            </Field>
          )}
        />
        {requirePasswordConfirm && (
          <form.Field
            name="passwordConfirm"
            children={(field) => (
              <Field id={field.name}>
                <FieldLabel>Confirm your password</FieldLabel>
                <LMInputField
                  name={field.name}
                  type="password"
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  autoComplete="new-password"
                  isInvalid={isFieldInvalid(field)}
                  errorMessage={getFirstError(field)}
                />
              </Field>
            )}
          />
        )}
      </FieldGroup>
    </form>
  );
}
