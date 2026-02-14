import { createFileRoute } from "@tanstack/react-router";
import { MainPanelContainer } from "@/components/leszmonitor/MainPanelContainer.tsx";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  TypographyH1,
  TypographyH2,
} from "@/components/leszmonitor/sidebar/Typography.tsx";
import { addGroup, getGroups, type GroupInput } from "@/lib/data/groupData.ts";
import { Button } from "@/components/ui/button.tsx";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card.tsx";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field.tsx";
import { Input } from "@/components/ui/input.tsx";
import { useForm } from "@tanstack/react-form";
import { z } from "zod";
import { Textarea } from "@/components/ui/textarea.tsx";
import { type ColumnDef, getCoreRowModel } from "@tanstack/table-core";
import type { Group } from "@/lib/types.ts";
import { flexRender, useReactTable } from "@tanstack/react-table";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx";
import { StyledLink } from "@/components/leszmonitor/StyledLink.tsx";
import { formatDate } from "@/lib/utils.ts";

export const Route = createFileRoute("/_authenticated/team/$teamId/groups/")({
  component: Groups,
});

const groupFormSchema = z.object({
  name: z.string().min(1, "Group name is required"),
  displayId: z.string().min(1, "Display ID is required"),
  description: z.string(),
});

function Groups() {
  const teamId = Route.useParams().teamId;
  const queryClient = useQueryClient();

  const { data } = useQuery({
    queryKey: ["groups", teamId],
    queryFn: () => getGroups(teamId),
  });

  const addGroupMutation = useMutation({
    mutationFn: (newGroup: GroupInput) => addGroup(teamId, newGroup),
    onSuccess: () => {
      form.reset();
      queryClient.invalidateQueries({ queryKey: ["groups", teamId] });
    },
  });

  const form = useForm({
    defaultValues: {
      name: "",
      displayId: "",
      description: "",
    },
    validators: {
      onSubmit: groupFormSchema,
    },
    onSubmit: async ({ value }) => {
      await addGroupMutation.mutateAsync(value);
    },
  });

  const columns: ColumnDef<Group>[] = [
    {
      accessorKey: "name",
      header: "Name",
      cell: ({ row }) => {
        const name = row.original.name;
        return (
          <StyledLink
            to="/team/$teamId/groups/$groupId"
            params={{ teamId, groupId: row.original.displayId }}
          >
            {name}
          </StyledLink>
        );
      },
    },
    {
      accessorKey: "displayId",
      header: "Display ID",
    },
    {
      accessorKey: "createdAt",
      header: "Created At",
      cell: ({ row }) => {
        return formatDate(row.original.createdAt);
      },
    },
    {
      accessorKey: "updatedAt",
      header: "Updated At",
      cell: ({ row }) => {
        return formatDate(row.original.updatedAt);
      },
    },
    {
      accessorKey: "description",
      header: "Description",
    },
  ];

  const table = useReactTable({
    data: data || [],
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  if (!data) {
    return null;
  }

  return (
    <MainPanelContainer>
      <TypographyH1>Groups</TypographyH1>
      <Card>
        <CardHeader>
          <TypographyH2>Create New Group</TypographyH2>
        </CardHeader>
        <CardContent>
          <form
            id="group-form"
            onSubmit={(e) => {
              e.preventDefault();
              form.handleSubmit();
            }}
          >
            <FieldGroup className="gap-2">
              <div className="flex gap-8">
                <form.Field
                  name="name"
                  children={(field) => {
                    const isInvalid =
                      field.state.meta.isTouched && !field.state.meta.isValid;
                    return (
                      <Field>
                        <FieldLabel htmlFor={field.name}>Group Name</FieldLabel>
                        <Input
                          id={field.name}
                          name={field.name}
                          value={field.state.value}
                          onChange={(e) => field.handleChange(e.target.value)}
                          autoComplete="off"
                        />
                        {isInvalid && (
                          <FieldError errors={field.state.meta.errors} />
                        )}
                      </Field>
                    );
                  }}
                />
                <form.Field
                  name="displayId"
                  children={(field) => {
                    const isInvalid =
                      field.state.meta.isTouched && !field.state.meta.isValid;
                    return (
                      <Field>
                        <FieldLabel htmlFor={field.name}>Display ID</FieldLabel>
                        <Input
                          id={field.name}
                          name={field.name}
                          value={field.state.value}
                          onChange={(e) => field.handleChange(e.target.value)}
                          autoComplete="off"
                        />
                        {isInvalid && (
                          <FieldError errors={field.state.meta.errors} />
                        )}
                      </Field>
                    );
                  }}
                />
              </div>
              <form.Field
                name="description"
                children={(field) => {
                  return (
                    <Field>
                      <FieldLabel htmlFor={field.name}>
                        Description (Optional)
                      </FieldLabel>
                      <Textarea
                        id={field.name}
                        name={field.name}
                        value={field.state.value}
                        onChange={(e) => field.handleChange(e.target.value)}
                        autoComplete="off"
                      />
                    </Field>
                  );
                }}
              />
            </FieldGroup>
          </form>
        </CardContent>
        <CardFooter className="justify-end">
          <Button
            type="submit"
            form="group-form"
            disabled={addGroupMutation.isPending}
          >
            {addGroupMutation.isPending ? "Adding..." : "Add Group"}
          </Button>
        </CardFooter>
      </Card>
      <Card>
        <CardHeader>
          <TypographyH2>Existing Groups</TypographyH2>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              {table.getHeaderGroups().map((headerGroup) => (
                <TableRow key={headerGroup.id}>
                  {headerGroup.headers.map((header) => {
                    return (
                      <TableHead key={header.id}>
                        {header.isPlaceholder
                          ? null
                          : flexRender(
                              header.column.columnDef.header,
                              header.getContext(),
                            )}
                      </TableHead>
                    );
                  })}
                </TableRow>
              ))}
            </TableHeader>
            <TableBody>
              {table.getRowModel().rows?.length ? (
                table.getRowModel().rows.map((row) => (
                  <TableRow
                    key={row.id}
                    data-state={row.getIsSelected() && "selected"}
                  >
                    {row.getVisibleCells().map((cell) => (
                      <TableCell key={cell.id}>
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext(),
                        )}
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell
                    colSpan={columns.length}
                    className="h-24 text-center"
                  >
                    No results.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </MainPanelContainer>
  );
}
