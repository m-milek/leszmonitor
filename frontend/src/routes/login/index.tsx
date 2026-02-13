import { createFileRoute, useNavigate } from "@tanstack/react-router";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card.tsx";
import { Logo } from "@/components/leszmonitor/Logo.tsx";
import { Button } from "@/components/ui/button.tsx";
import { z } from "zod";
import { useForm } from "@tanstack/react-form";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field.tsx";
import { Input } from "@/components/ui/input.tsx";
import { fetchLoginToken } from "@/lib/fetchLoginToken.ts";
import { userAtom, usernameAtom } from "@/lib/atoms.ts";
import { useSetAtom } from "jotai";
import { jwtDecode } from "jwt-decode";
import type { JwtClaims } from "@/lib/types.ts";
import { fetchUser } from "@/lib/data/userData.ts";

export const Route = createFileRoute("/login/")({
  component: RouteComponent,
});
const loginFormSchema = z.object({
  username: z.string().min(1, "Username is required"),
  password: z.string().min(1, "Password is required"),
});

function RouteComponent() {
  const navigate = useNavigate();

  const setUsername = useSetAtom(usernameAtom);
  const setUser = useSetAtom(userAtom);

  const form = useForm({
    defaultValues: {
      username: "",
      password: "",
    },
    validators: {
      onSubmit: loginFormSchema,
    },
    onSubmit: async ({ value }) => {
      const loginResponse = await fetchLoginToken(value);
      await cookieStore.set({
        name: "LOGIN_TOKEN",
        value: loginResponse.jwt,
      });

      const claims = jwtDecode(loginResponse.jwt) as JwtClaims;
      setUsername(claims.username);

      const user = await fetchUser(claims.username);
      setUser(user);

      await navigate({ to: "/", replace: true });
    },
  });

  return (
    <main className="h-screen w-screen">
      <div className="flex h-full w-full items-center justify-center">
        <Card className="w-full max-w-sm">
          <CardHeader className="text-center">
            <CardTitle className="flex flex-col items-center">
              <Logo />
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
                  children={(field) => {
                    const isInvalid =
                      field.state.meta.isTouched && !field.state.meta.isValid;
                    return (
                      <Field>
                        <FieldLabel htmlFor={field.name}>Username</FieldLabel>
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
                  name="password"
                  children={(field) => {
                    const isInvalid =
                      field.state.meta.isTouched && !field.state.meta.isValid;
                    return (
                      <Field>
                        <FieldLabel htmlFor={field.name}>Password</FieldLabel>
                        <Input
                          id={field.name}
                          name={field.name}
                          type="password"
                          value={field.state.value}
                          onChange={(e) => field.handleChange(e.target.value)}
                          autoComplete="current-password"
                        />
                        {isInvalid && (
                          <FieldError errors={field.state.meta.errors} />
                        )}
                      </Field>
                    );
                  }}
                />
              </FieldGroup>
            </form>
          </CardContent>
          <CardFooter>
            <Button
              className="w-full"
              onClick={() => form.handleSubmit()}
              type="submit"
              form="login-form"
            >
              Log in
            </Button>
          </CardFooter>
        </Card>
      </div>
    </main>
  );
}
