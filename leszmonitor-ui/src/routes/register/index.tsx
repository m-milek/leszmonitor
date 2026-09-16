import { createFileRoute } from "@tanstack/react-router";
import { RegisterPage } from "@/features/auth/pages/RegisterPage.tsx";

export const Route = createFileRoute("/register/")({
  head: () => ({
    meta: [{ title: "Register | Leszmonitor" }],
  }),
  component: RegisterPage,
});
