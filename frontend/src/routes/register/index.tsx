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
import { getUser, registerUser } from "@/lib/data/userData.ts";
import { isJwtClaims } from "@/lib/jwt.ts";
import { setCookie } from "@/lib/cookies.ts";
import { toast } from "sonner";
import { RegisterUserForm } from "@/components/leszmonitor/forms/RegisterUserForm.tsx";
import { type RegisterUserPayload } from "@/lib/data/userData.ts";

export const Route = createFileRoute("/register/")({
  component: RegisterComponent,
});



function RegisterComponent() {
  const navigate = useNavigate();

  const { setUsername, setUser } = useAppStore();

  const handleSubmit = async (value: RegisterUserPayload) => {
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

      const user = await getUser(claims.username);
      setUser(user);

      await navigate({ to: "/", replace: true });
    } catch (error) {
      console.error("Registration failed:", error);
      toast.error("Registration failed. Please try again.");
    }
  };

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
            <RegisterUserForm id="login-form" onSubmit={handleSubmit} />
          </CardContent>
          <CardFooter>
            <Button
              className="w-full"
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
