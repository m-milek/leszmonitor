import { useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { LeszmonitorLogo } from "@/components/common/LeszmonitorLogo.tsx";
import { RegisterUserForm } from "@/features/auth/forms/RegisterUserForm.tsx";
import {
  registerUser,
  type RegisterUserPayload,
} from "@/features/users/users-api.ts";
import { establishSession } from "@/features/auth/lib/session.ts";
import { useAppStore } from "@/app/store.ts";

export function RegisterPage() {
  const navigate = useNavigate();

  const { setUsername, setUser } = useAppStore();

  const handleSubmit = async (value: RegisterUserPayload) => {
    try {
      console.log("Registering user with values:", value);
      await registerUser(value);

      const established = await establishSession(value, {
        setUsername,
        setUser,
      });
      if (!established) {
        return;
      }

      await navigate({ to: "/monitors", replace: true });
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
            <Button className="w-full" type="submit" form="login-form">
              Register
            </Button>
          </CardFooter>
        </Card>
      </div>
    </main>
  );
}
