import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export const RESULTS_RANGE_HOURS = {
  "1h": 1,
  "4h": 4,
  "12h": 12,
  "24h": 24,
} as const;

export type ResultsRange = keyof typeof RESULTS_RANGE_HOURS;

export const resultsRangeStart = (range: ResultsRange) =>
  new Date(Date.now() - RESULTS_RANGE_HOURS[range] * 60 * 60 * 1000);

export const resultsRangePerPage = (
  range: ResultsRange,
  intervalSeconds: number,
) => Math.ceil((RESULTS_RANGE_HOURS[range] * 60 * 60) / intervalSeconds) * 2;

export interface ResultsRangeSelectProps {
  value: ResultsRange;
  onChange: (value: ResultsRange) => void;
}

export const ResultsRangeSelect = ({
  value,
  onChange,
}: ResultsRangeSelectProps) => (
  <Select
    value={value}
    onValueChange={(next) => next && onChange(next as ResultsRange)}
  >
    <SelectTrigger>
      <SelectValue />
    </SelectTrigger>
    <SelectContent alignItemWithTrigger={false}>
      {Object.keys(RESULTS_RANGE_HOURS).map((range) => (
        <SelectItem key={range} value={range}>
          Last {range}
        </SelectItem>
      ))}
    </SelectContent>
  </Select>
);
