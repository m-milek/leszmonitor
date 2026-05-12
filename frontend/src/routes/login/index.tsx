import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card.tsx";
import { LeszmonitorLogo } from "@/components/leszmonitor/ui/LeszmonitorLogo.tsx";
import { Button } from "@/components/ui/button.tsx";
import { z } from "zod";
import { useForm } from "@tanstack/react-form";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field.tsx";
import { fetchLoginToken } from "@/lib/fetchLoginToken.ts";
import { useAppStore } from "@/lib/store.ts";
import { jwtDecode } from "jwt-decode";
import { isJwtClaims } from "@/lib/jwt.ts";
import { fetchUser } from "@/lib/data/userData.ts";
import { setCookie } from "@/lib/cookies.ts";
import { LMInputField } from "@/components/leszmonitor/forms/inputs/LMInputField.tsx";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/leszmonitor/forms/inputs/utils.ts";
import { toast } from "sonner";

export const Route = createFileRoute("/login/")({
  component: RouteComponent,
});

const loginFormSchema = z.object({
  username: z.string().min(1, "Username is required"),
  password: z.string().min(1, "Password is required"),
});

function RouteComponent() {
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
        const loginResponse = await fetchLoginToken(value);

        setCookie("LOGIN_TOKEN", loginResponse.jwt, {
          maxAge: 24 * 60 * 60,
          path: "/",
          sameSite: "Lax",
        });

        const claims = jwtDecode(loginResponse.jwt);
        if (!isJwtClaims(claims)) {
          console.error("Invalid JWT claims");
          return;
        }

        setUsername(claims.username);

        const user = await fetchUser(claims.username);
        setUser(user);

        await navigate({ to: "/", replace: true });
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
                    <Field>
                      <FieldLabel htmlFor={field.name}>Username</FieldLabel>
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
                    <Field>
                      <FieldLabel htmlFor={field.name}>Password</FieldLabel>
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
            <Button
              className="w-full"
              onClick={() => form.handleSubmit()}
              type="submit"
              form="login-form"
            >
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
