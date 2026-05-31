import { ErrorTooltip } from "@/components/leszmonitor/forms/inputs/ErrorTooltip.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Plus, X } from "lucide-react";
import { memo, useCallback, useState } from "react";

interface ListItem {
  id: string;
  value: string;
}

export interface LMListInputProps {
  name: string;
  value?: string[];
  onChange: (value: string[]) => void;
  placeholder?: string;
  addButtonText?: string;
  isInvalid?: boolean;
  errorMessage?: string;
  disabled?: boolean;
}

let nextId = 0;

function normalizeValues(values: ListItem[]): string[] {
  return values
    .map((item) => item.value.trim())
    .filter((value) => value.length > 0);
}

interface ListRowProps {
  item: ListItem;
  placeholder?: string;
  onChange: (id: string, value: string) => void;
  onDelete: (id: string) => void;
  disabled?: boolean;
}

const ListRow = memo(function ListRow({
  item,
  placeholder,
  onChange,
  onDelete,
  disabled,
}: ListRowProps) {
  return (
    <div className="flex gap-2">
      <Input
        type="text"
        value={item.value}
        onChange={(e) => onChange(item.id, e.target.value)}
        placeholder={placeholder ?? "Value"}
        className="flex-1"
        disabled={disabled}
      />
      <Button
        type="button"
        aria-label="Remove"
        aria-roledescription={`Remove item ${item.value}`}
        onClick={() => onDelete(item.id)}
        variant="ghost"
        size="sm"
        className="px-2"
        disabled={disabled}
      >
        <X className="h-4 w-4" />
      </Button>
    </div>
  );
});

export function LMListInput(props: Readonly<LMListInputProps>) {
  const [items, setItems] = useState<ListItem[]>(() =>
    (props.value ?? []).map((value) => ({
      id: `li-${nextId++}`,
      value,
    })),
  );

  const update = useCallback(
    (updated: ListItem[]) => {
      setItems(updated);
      props.onChange(normalizeValues(updated));
    },
    [props],
  );

  const handleChange = useCallback(
    (id: string, value: string) => {
      update(items.map((item) => (item.id === id ? { ...item, value } : item)));
    },
    [items, update],
  );

  const handleDelete = useCallback(
    (id: string) => {
      update(items.filter((item) => item.id !== id));
    },
    [items, update],
  );

  const handleAdd = useCallback(() => {
    setItems((prev) => [...prev, { id: `li-${nextId++}`, value: "" }]);
  }, []);

  return (
    <ErrorTooltip isOpen={props.isInvalid} message={props.errorMessage}>
      <div className="flex flex-col gap-2">
        {items.map((item) => (
          <ListRow
            key={item.id}
            item={item}
            placeholder={props.placeholder}
            onChange={handleChange}
            onDelete={handleDelete}
            disabled={props.disabled}
          />
        ))}
        <Button
          type="button"
          onClick={handleAdd}
          variant="ghost"
          size="sm"
          disabled={props.disabled}
        >
          <Plus className="h-4 w-4" />
          {props.addButtonText ?? "Add"}
        </Button>
      </div>
    </ErrorTooltip>
  );
}
