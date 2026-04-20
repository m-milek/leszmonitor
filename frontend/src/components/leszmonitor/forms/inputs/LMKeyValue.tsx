import { ErrorTooltip } from "@/components/leszmonitor/forms/inputs/ErrorTooltip.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useState, useCallback, memo } from "react";
import { Plus, X } from "lucide-react";

interface KeyValuePair {
  id: string;
  key: string;
  value: string;
}

export interface LMKeyValueInputProps {
  value?: Record<string, string>;
  onChange: (value: Record<string, string>) => void;
  name: string;
  isInvalid?: boolean;
  errorMessage?: string;
  keyPlaceholder?: string;
  valuePlaceholder?: string;
  addButtonText?: string;
}

let nextId = 0;

function pairsToRecord(pairs: KeyValuePair[]): Record<string, string> {
  return Object.fromEntries(
    pairs.filter((p) => p.key || p.value).map((p) => [p.key, p.value]),
  );
}

interface KeyValueRowProps {
  pair: KeyValuePair;
  keyPlaceholder?: string;
  valuePlaceholder?: string;
  onKeyChange: (id: string, value: string) => void;
  onValueChange: (id: string, value: string) => void;
  onDelete: (id: string) => void;
}

const KeyValueRow = memo(function KeyValueRow({
  pair,
  keyPlaceholder,
  valuePlaceholder,
  onKeyChange,
  onValueChange,
  onDelete,
}: KeyValueRowProps) {
  return (
    <div className="flex gap-2">
      <Input
        type="text"
        value={pair.key}
        onChange={(e) => onKeyChange(pair.id, e.target.value)}
        placeholder={keyPlaceholder ?? "Key"}
        className="flex-1"
      />
      <Input
        type="text"
        value={pair.value}
        onChange={(e) => onValueChange(pair.id, e.target.value)}
        placeholder={valuePlaceholder ?? "Value"}
        className="flex-1"
      />
      <Button
        type="button"
        onClick={() => onDelete(pair.id)}
        variant="ghost"
        size="sm"
        className="px-2"
      >
        <X className="h-4 w-4" />
      </Button>
    </div>
  );
});

export function LMKeyValueInput(props: LMKeyValueInputProps) {
  const [pairs, setPairs] = useState<KeyValuePair[]>(() =>
    Object.entries(props.value ?? {}).map(([key, value]) => ({
      id: `kv-${nextId++}`,
      key,
      value,
    })),
  );

  const update = useCallback(
    (updated: KeyValuePair[]) => {
      setPairs(updated);
      props.onChange(pairsToRecord(updated));
    },
    [props],
  );

  const handleKeyChange = useCallback(
    (id: string, val: string) => {
      update(pairs.map((p) => (p.id === id ? { ...p, key: val } : p)));
    },
    [pairs, update],
  );

  const handleValueChange = useCallback(
    (id: string, val: string) => {
      update(pairs.map((p) => (p.id === id ? { ...p, value: val } : p)));
    },
    [pairs, update],
  );

  const handleDelete = useCallback(
    (id: string) => {
      update(pairs.filter((p) => p.id !== id));
    },
    [pairs, update],
  );

  const handleAdd = useCallback(() => {
    setPairs((prev) => [...prev, { id: `kv-${nextId++}`, key: "", value: "" }]);
  }, []);

  return (
    <ErrorTooltip isOpen={props.isInvalid} message={props.errorMessage}>
      <div className="flex flex-col gap-2">
        {pairs.map((pair) => (
          <KeyValueRow
            key={pair.id}
            pair={pair}
            keyPlaceholder={props.keyPlaceholder}
            valuePlaceholder={props.valuePlaceholder}
            onKeyChange={handleKeyChange}
            onValueChange={handleValueChange}
            onDelete={handleDelete}
          />
        ))}
        <Button type="button" onClick={handleAdd} variant="ghost" size="sm">
          <Plus className="h-4 w-4" />
          {props.addButtonText ?? "Add"}
        </Button>
      </div>
    </ErrorTooltip>
  );
}
