import { createFileRoute } from "@tanstack/react-router";
import { PageContainer } from "@/components/leszmonitor/PageContainer.tsx";
import {
  TypographyH1,
  TypographyH2,
} from "@/components/leszmonitor/ui/Typography.tsx";
import { useQuery } from "@tanstack/react-query";
import { getAllUsers } from "@/lib/data/userData.ts";
import { UsersTable } from "@/components/leszmonitor/tables/UsersTable.tsx";
import { Skeleton } from "@/components/ui/skeleton.tsx";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx";
import { Button } from "@/components/ui/button.tsx";
import { RegisterUserForm } from "@/components/leszmonitor/forms/RegisterUserForm.tsx";
import { registerUser, type RegisterUserPayload } from "@/lib/data/userData.ts";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { Plus } from "lucide-react";
import { useState } from "react";

export const Route = createFileRoute("/_authenticated/admin/")({
  component: AdminDashboardRoute,
});

function AdminDashboardRoute() {
  const queryClient = useQueryClient();
  const { data: users } = useQuery({
    queryKey: ["users"],
    queryFn: () => getAllUsers(),
  });

  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const registerMutation = useMutation({
    mutationFn: (values: RegisterUserPayload) => registerUser(values),
    onSuccess: () => {
      toast.success("User registered successfully");
      setIsDialogOpen(false);
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
    onError: (error) => {
      toast.error("Failed to register user: " + error.message);
    },
  });

  return (
    <PageContainer>
      <TypographyH1>Administration Dashboard</TypographyH1>
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
          <div className="space-y-1">
            <TypographyH2>Users</TypographyH2>
            <p className="text-sm text-muted-foreground">
              All users in this Leszmonitor instance
            </p>
          </div>
          <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
              <Button size="sm">
                <Plus className="mr-2 h-4 w-4" />
                Add User
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add New User</DialogTitle>
                <DialogDescription>
                  Register a new user in the system.
                </DialogDescription>
              </DialogHeader>
              <RegisterUserForm
                id="add-user-form"
                requirePasswordConfirm={false}
                onSubmit={async (values) => {
                  await registerMutation.mutateAsync(values);
                }}
              />
              <Button type="submit" form="add-user-form" className="w-full">
                Create User
              </Button>
            </DialogContent>
          </Dialog>
        </CardHeader>
        <CardContent>
          {users ? <UsersTable users={users} /> : <Skeleton className="h-24" />}
        </CardContent>
      </Card>
    </PageContainer>
  );
}
