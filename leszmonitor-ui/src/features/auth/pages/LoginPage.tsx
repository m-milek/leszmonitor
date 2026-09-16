import { Link, useNavigate } from "@tanstack/react-router";
import { z } from "zod";
import { useForm } from "@tanstack/react-form";
import { toast } from "sonner";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { LeszmonitorLogo } from "@/components/common/LeszmonitorLogo";
import { LMInputField } from "@/components/form/LMInputField";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/form/field-state";
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
          toast.error(
            "Failed to log in. Please check your credentials and try again.",
          );
        }
      }
    },
  });

  return (
    <main className="h-screen w-screen bg-background">
      <div className="flex h-full w-full items-center justify-center">
        <Card className="w-full max-w-sm">
          <CardHeader className="text-center">
            <CardTitle className="flex flex-col items-center">
              <LeszmonitorLogo />
            </CardTitle>
            <CardDescription>Log in to Leszmonitor</CardDescription>
          </CardHeader>
          <CardContent>
            <form
              id="login-form"
              onSubmit={(e) => {
                e.preventDefault();
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
                        autoComplete="current-password"
                        isInvalid={isFieldInvalid(field)}
                        errorMessage={getFirstError(field)}
                      />
                    </Field>
                  )}
                />
              </FieldGroup>
            </form>
          </CardContent>
          <CardFooter className="flex-col items-center gap-4">
            <Button className="w-full" type="submit" form="login-form">
              Log in
            </Button>
            <Link to="/register" className="text-sm text-primary">
              Don&#39;t have an account?
            </Link>
          </CardFooter>
        </Card>
      </div>
    </main>
  );
}
