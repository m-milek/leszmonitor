import type { ReactNode } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { LeszmonitorLogo } from "@/components/common/LeszmonitorLogo";

export interface AuthCardLayoutProps {
  description: ReactNode;
  children: ReactNode;
  footer: ReactNode;
}

export const AuthCardLayout = ({
  description,
  children,
  footer,
}: AuthCardLayoutProps) => (
  <main className="flex h-svh w-full items-center justify-center">
    <Card className="w-full max-w-sm">
      <CardHeader className="text-center">
        <CardTitle className="flex flex-col items-center">
          <LeszmonitorLogo />
        </CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>{children}</CardContent>
      <CardFooter className="flex-col items-center gap-4">{footer}</CardFooter>
    </Card>
  </main>
);
