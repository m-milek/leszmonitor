import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { toast } from "@/components/ui/toast";
import {
  GlobalParametersApi,
  type GlobalParameterValue,
} from "@/features/instance/global-parameters-api";
import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { GlobalParametersForm } from "@/features/instance/forms/GlobalParametersForm";

const FORM_ID = "global-parameters-form";

export const GlobalParametersConfiguration = () => {
  const queryClient = useQueryClient();
  const { data: globalParameters, dataUpdatedAt } = useQuery({
    queryKey: ["globalParameters"],
    queryFn: () => GlobalParametersApi.getAll(),
  });

  const [hasChanges, setHasChanges] = useState(false);

  const updateMutation = useMutation({
    mutationFn: (values: Record<string, GlobalParameterValue | null>) =>
      GlobalParametersApi.update(values),
    onSuccess: () => {
      toast.add({ title: "Global parameters saved", type: "success" });
      queryClient.invalidateQueries({ queryKey: ["globalParameters"] });
    },
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle>Global Parameters</CardTitle>
      </CardHeader>
      <CardContent>
        {globalParameters && (
          <GlobalParametersForm
            key={dataUpdatedAt}
            id={FORM_ID}
            globalParameters={globalParameters}
            onHasChangesChange={setHasChanges}
            onSubmit={async (values) => {
              await updateMutation.mutateAsync(values);
            }}
          />
        )}
      </CardContent>
      <CardFooter>
        <Button
          type="submit"
          form={FORM_ID}
          disabled={!hasChanges || updateMutation.isPending}
        >
          {updateMutation.isPending ? "Saving..." : "Save"}
        </Button>
      </CardFooter>
    </Card>
  );
};
