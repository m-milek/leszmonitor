import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { TypographyH1 } from "@/components/leszmonitor/Typography.tsx";
import { useForm } from "@tanstack/react-form";
import {
  newMonitorSchema,
  newMonitorSchemaDefaultValues,
} from "@/lib/types.ts";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card.tsx";
import { Button } from "@/components/ui/button.tsx";
import { LMInputField } from "@/components/leszmonitor/forms/LMInputField.tsx";
import * as React from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx";
import { Field, FieldError, FieldLabel } from "@/components/ui/field.tsx";

export const Route = createFileRoute(
  "/_authenticated/team/$teamId/monitors/new/",
)({
  component: NewMonitorComponent,
});

function NewMonitorComponent() {
  const form = useForm({
    defaultValues: newMonitorSchemaDefaultValues,
    validators: {
      onSubmit: newMonitorSchema,
    },
    onSubmit: ({ value }) => {
      console.log(value);
    },
    onSubmitInvalid: ({ value }) => {
      console.log("Invalid form submission");
      console.log("Values:", value);
    },
  });

  const onSubmit = (e: React.SubmitEvent) => {
    console.log("Submitting form");
    e.preventDefault();
    form.handleSubmit();
  };

  return (
    <MainPanelContainer>
      <TypographyH1>New Monitor Wizard</TypographyH1>
      <Card>
        <CardHeader>Form</CardHeader>
        <CardContent>
          <form
            id="new-monitor-form"
            onSubmit={onSubmit}
            className="flex items-end gap-4"
          >
            <form.Field
              name="name"
              children={(field) => {
                return <LMInputField label="Name" field={field} />;
              }}
            />
            <form.Field
              name="displayId"
              children={(field) => {
                return <LMInputField label="Display ID" field={field} />;
              }}
            />
            <form.Field
              name="interval"
              children={(field) => {
                return <LMInputField label="Interval (s)" field={field} />;
              }}
            />
            <form.Field
              name={"type"}
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;
                return (
                  <Field>
                    <FieldLabel>Type</FieldLabel>
                    <Select>
                      <SelectTrigger>
                        <SelectValue placeholder="Select Monitor Type" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value={"http"}>HTTP</SelectItem>
                        <SelectItem value={"ping"}>Ping</SelectItem>
                      </SelectContent>
                    </Select>
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                );
              }}
            />
            {form.state.values.type === "http" && <div>Dupa</div>}
            {form.state.values.type === "ping" && <div>Ping</div>}
          </form>
        </CardContent>
        <CardFooter>
          <Button type="submit" form="new-monitor-form">
            Create Monitor
          </Button>
        </CardFooter>
      </Card>
    </MainPanelContainer>
  );
}
