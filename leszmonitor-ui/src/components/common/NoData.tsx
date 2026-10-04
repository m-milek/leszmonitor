import { Minus } from "lucide-react";

export interface NoDataProps {
  label?: string;
}

export const NoData = ({ label = "No data" }: NoDataProps) => (
  <Minus className="size-4 text-muted-foreground" aria-label={label} />
);
