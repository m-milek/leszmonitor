import { Fragment, useEffect } from "react";
import { useForm, useStore } from "@tanstack/react-form";
import {
  FieldGroup,
  FieldLegend,
  FieldSeparator,
  FieldSet,
} from "@/components/ui/field";
import { GlobalParameterRow } from "@/features/instance/components/GlobalParameterRow";
import type {
  GlobalParameter,
  GlobalParameterValue,
} from "@/features/instance/global-parameters-api";

export interface GlobalParametersFormProps {
  id?: string;
  globalParameters: GlobalParameter[];
  onSubmit: (
    values: Record<string, GlobalParameterValue | null>,
  ) => Promise<void>;
  onHasChangesChange?: (hasChanges: boolean) => void;
}

const sectionName = (key: string) => {
  const prefix = key.split(".")[0];
  return prefix.charAt(0).toUpperCase() + prefix.slice(1);
};

export function GlobalParametersForm({
  id = "global-parameters-form",
  globalParameters,
  onSubmit,
  onHasChangesChange,
}: Readonly<GlobalParametersFormProps>) {
  const form = useForm({
    defaultValues: {
      values: globalParameters.map(
        (param): GlobalParameterValue | null => param.value,
      ),
    },
    onSubmit: async ({ value }) => {
      const changed: Record<string, GlobalParameterValue | null> = {};
      globalParameters.forEach((param, index) => {
        if (value.values[index] !== param.value) {
          changed[param.key] = value.values[index];
        }
      });

      if (Object.keys(changed).length > 0) {
        await onSubmit(changed);
      }
    },
  });

  const hasChanges = useStore(form.store, (state) => !state.isDefaultValue);
  useEffect(() => {
    onHasChangesChange?.(hasChanges);
  }, [hasChanges, onHasChangesChange]);

  const sections = Map.groupBy(
    globalParameters.map((param, index) => ({ param, index })),
    ({ param }) => sectionName(param.key),
  );

  return (
    <form
      id={id}
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}
    >
      <FieldGroup>
        {[...sections].map(([section, entries], sectionIndex) => (
          <Fragment key={section}>
            {sectionIndex > 0 && <FieldSeparator />}
            <FieldSet>
              <FieldLegend>{section}</FieldLegend>
              {entries.map(({ param, index }) => (
                <form.Field key={param.key} name={`values[${index}]`}>
                  {(field) => (
                    <GlobalParameterRow
                      id={field.name}
                      parameter={param}
                      value={field.state.value}
                      onChange={field.handleChange}
                    />
                  )}
                </form.Field>
              ))}
            </FieldSet>
          </Fragment>
        ))}
      </FieldGroup>
    </form>
  );
}
