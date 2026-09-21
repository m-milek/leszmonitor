import { Link, useNavigate } from "@tanstack/react-router";
import { z } from "zod";
import { useForm } from "@tanstack/react-form";
import { toast } from "@/components/ui/toast";
import { Button } from "@/components/ui/button";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { AuthCardLayout } from "@/features/auth/components/AuthCardLayout";
import { LMInputField } from "@/components/form/LMInputField";
import { getFirstError, isFieldInvalid } from "@/components/form/field-state";
import { establishSession } from "@/features/auth/lib/session";
import { useAppStore } from "@/app/store";

const loginFormSchema = z.object({
  username: z.string().min(1, "Username is required"),
  password: z.string().min(1, "Password is required"),
});

export function LoginPage() {
  const navigate = useNavigate();

  const { setUsername, setUser } = useAppStore();

  const form = useForm({
    defaultValues: {
      username: "",
      password: "",
    },
    validators: {
      onSubmit: loginFormSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        const established = await establishSession(value, {
          setUsername,
          setUser,
        });
        if (!established) {
          return;
        }

        await navigate({ to: "/monitors", replace: true });
      } catch (error) {
        if (error instanceof Error) {
          console.error(error);
          toast.add({
            title:
              "Failed to log in. Please check your credentials and try again.",
            type: "error",
          });
        }
      }
    },
  });

  return (
    <AuthCardLayout
      description="Log in to Leszmonitor"
      footer={
        <>
          <Button className="w-full" type="submit" form="login-form">
            Log in
          </Button>
          <Link to="/register" className="text-sm text-primary">
            Don&#39;t have an account?
          </Link>
        </>
      }
    >
      <form
        id="login-form"
        onSubmit={(e) => {
          e.preventDefault();
          form.handleSubmit();
        }}
      >
        <FieldGroup>
          <form.Field name="username">
            {(field) => (
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
          </form.Field>
          <form.Field name="password">
            {(field) => (
              <Field id={field.name}>
                <FieldLabel>Password</FieldLabel>
                <LMInputField
                  name={field.name}
                  type="password"
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                  autoComplete="current-password"
                  isInvalid={isFieldInvalid(field)}
                  errorMessage={getFirstError(field)}
                />
              </Field>
            )}
          </form.Field>
        </FieldGroup>
      </form>
    </AuthCardLayout>
  );
}
