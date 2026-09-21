import { UsersApi } from "@/features/users/users-api";
import { type RegisterUserPayload } from "@/features/users/users-api";
import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Plus } from "lucide-react";
import { toast } from "@/components/ui/toast";
import { PageContainer } from "@/components/common/PageContainer";
import { TypographyH1 } from "@/components/common/Typography";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { UsersTable } from "@/features/users/components/UsersTable";
import { RegisterUserForm } from "@/features/auth/forms/RegisterUserForm";

export function AdminPage() {
  const queryClient = useQueryClient();
  const { data: users } = useQuery({
    queryKey: ["users"],
    queryFn: () => UsersApi.getAll(),
  });

  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const registerMutation = useMutation({
    mutationFn: (values: RegisterUserPayload) => UsersApi.register(values),
    onSuccess: () => {
      toast.add({ title: "User registered successfully", type: "success" });
      setIsDialogOpen(false);
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
    onError: (error) => {
      toast.add({
        title: "Failed to register user: " + error.message,
        type: "error",
      });
    },
  });

  return (
    <PageContainer>
      <TypographyH1>Administration Dashboard</TypographyH1>
      <Card>
        <CardHeader>
          <CardTitle>Users</CardTitle>
          <CardDescription>
            All users in this Leszmonitor instance
          </CardDescription>
          <CardAction>
            <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
              <DialogTrigger render={<Button size="sm" />}>
                <Plus />
                Add User
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
          </CardAction>
        </CardHeader>
        <CardContent className="px-0">
          {users && <UsersTable users={users} />}
        </CardContent>
      </Card>
    </PageContainer>
  );
}
