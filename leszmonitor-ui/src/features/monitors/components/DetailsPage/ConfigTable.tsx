import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";
import {
  type FieldConfig,
  formatValue,
} from "@/features/monitors/components/DetailsPage/config-fields";

export interface ConfigTableProps<T> {
  data: T;
  fields: FieldConfig<T>;
}

export const ConfigTable = <T,>({ data, fields }: ConfigTableProps<T>) => (
  <Table>
    <TableBody>
      {(Object.keys(fields) as (keyof T)[]).map((key) => {
        const field = fields[key];
        const value = data[key];
        if (!field || value === undefined || value === null) return null;
        return (
          <TableRow key={String(key)} className="border-0">
            <TableCell className="p-0 align-top font-medium">
              {field.label}
            </TableCell>
            <TableCell className="p-0 whitespace-pre-wrap break-all">
              {field.render ? field.render(value, data) : formatValue(value)}
            </TableCell>
          </TableRow>
        );
      })}
    </TableBody>
  </Table>
);
