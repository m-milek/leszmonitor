import { createFileRoute, useNavigate } from "@tanstack/react-router";
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
import { fetchUser, registerUser } from "@/lib/data/userData.ts";
import { isJwtClaims } from "@/lib/jwt.ts";
import { setCookie } from "@/lib/cookies.ts";
import { toast } from "sonner";
import { LMInputField } from "@/components/leszmonitor/forms/inputs/LMInputField.tsx";
import {
  getFirstError,
  isFieldInvalid,
} from "@/components/leszmonitor/forms/inputs/utils.ts";

export const Route = createFileRoute("/register/")({
  component: RegisterComponent,
});

const registerFormSchema = z
  .object({
    username: z
      .string()
      .min(2, "Username has to be at least 2 characters long"),
    password: z
      .string()
      .min(6, "Password has to be at least 6 characters long"),
    passwordConfirm: z
      .string()
      .min(6, "Verify the password by entering it again"),
  })
  .refine((data) => data.password === data.passwordConfirm, {
    message: "Passwords don't match",
  });

function RegisterComponent() {
  const navigate = useNavigate();

  const { setUsername, setUser } = useAppStore();

  const form = useForm({
    defaultValues: {
      username: "",
      password: "",
      passwordConfirm: "",
    },
    validators: {
      onSubmit: registerFormSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        console.log("Registering user with values:", value);
        await registerUser(value);

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
        console.error("Registration failed:", error);
        toast.error("Registration failed. Please try again.");
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
            <CardDescription>
              Register a new account on Leszmonitor
            </CardDescription>
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
                        autoComplete="new-password"
                        isInvalid={isFieldInvalid(field)}
                        errorMessage={getFirstError(field)}
                      />
                    </Field>
                  )}
                />
                <form.Field
                  name="passwordConfirm"
                  children={(field) => (
                    <Field>
                      <FieldLabel htmlFor={field.name}>
                        Confirm your password
                      </FieldLabel>
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
              Register
            </Button>
          </CardFooter>
        </Card>
      </div>
    </main>
  );
}
