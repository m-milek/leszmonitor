type FieldLike = {
  state: {
    meta: {
      isTouched: boolean;
      isValid: boolean;
      errors?: Array<{ message?: string } | undefined>;
    };
  };
};

export const isFieldInvalid = (field: FieldLike) =>
  field.state.meta.isTouched && !field.state.meta.isValid;

export const getFirstError = (field: FieldLike) => {
  return field.state.meta.errors?.[0]?.message ?? "";
};
