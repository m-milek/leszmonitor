import { Field, FieldError, FieldLabel } from "@/components/ui/field.tsx";
import { Input } from "@/components/ui/input.tsx";

export interface LMFieldProps<T> {
  label: string;
  field: T;
}

export const LMInputField = <T,>({ label, field }: LMFieldProps<T>) => {
  //@ts-expect-error - we know this is a field
  const isInvalid = field.state.meta.isTouched && !field.state.meta.isValid;
  return (
    <Field>
      <FieldLabel>{label}</FieldLabel>
      <Input />
      {/*@ts-expect-error - we know this is a field */}
      {isInvalid && <FieldError errors={field.state.meta.errors} />}
    </Field>
  );
};
