import { UsersApi } from "@/features/users/users-api";
import { type RegisterUserPayload } from "@/features/users/users-api";
import { useNavigate } from "@tanstack/react-router";
import { toast } from "@/components/ui/toast";
import { Button } from "@/components/ui/button";
import { AuthCardLayout } from "@/features/auth/components/AuthCardLayout";
import { RegisterUserForm } from "@/features/auth/forms/RegisterUserForm";
import { establishSession } from "@/features/auth/lib/session";
import { useAppStore } from "@/app/store";

export function RegisterPage() {
  const navigate = useNavigate();

  const { setUsername, setUser } = useAppStore();

  const handleSubmit = async (value: RegisterUserPayload) => {
    try {
      console.log("Registering user with values:", value);
      await UsersApi.register(value);

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
      toast.add({
        title: "Registration failed. Please try again.",
        type: "error",
      });
    }
  };

  return (
    <AuthCardLayout
      description="Register a new account on Leszmonitor"
      footer={
        <Button className="w-full" type="submit" form="register-form">
          Register
        </Button>
      }
    >
      <RegisterUserForm id="register-form" onSubmit={handleSubmit} />
    </AuthCardLayout>
  );
}
