import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { formatDuration } from "@/lib/utils";

const HOUR = 60 * 60;

export const RESULTS_RANGE_OPTIONS = [HOUR, 4 * HOUR, 12 * HOUR, 24 * HOUR];

export const DEFAULT_RESULTS_RANGE = 24 * HOUR;

export const parseResultsRange = (value: unknown): number | undefined =>
  typeof value === "number" && Number.isInteger(value) && value > 0
    ? value
    : undefined;

export const resultsRangeStart = (rangeSeconds: number) =>
  new Date(Date.now() - rangeSeconds * 1000);

export const resultsRangePerPage = (
  rangeSeconds: number,
  intervalSeconds: number,
) => Math.ceil(rangeSeconds / intervalSeconds) * 2;

export interface ResultsRangeSelectProps {
  value: number;
  onChange: (value: number) => void;
}

export const ResultsRangeSelect = ({
  value,
  onChange,
}: ResultsRangeSelectProps) => {
  const options = RESULTS_RANGE_OPTIONS.includes(value)
    ? RESULTS_RANGE_OPTIONS
    : [...RESULTS_RANGE_OPTIONS, value].sort((a, b) => a - b);
  const items = options.map((range) => ({
    value: String(range),
    label: `Last ${formatDuration(range)}`,
  }));

  return (
    <Select
      items={items}
      value={String(value)}
      onValueChange={(next) => next && onChange(Number(next))}
    >
      <SelectTrigger>
        <SelectValue />
      </SelectTrigger>
      <SelectContent alignItemWithTrigger={false}>
        {items.map((item) => (
          <SelectItem key={item.value} value={item.value}>
            {item.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};
