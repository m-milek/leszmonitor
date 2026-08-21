import { cn } from "@/lib/utils.ts";

export const ShortId = ({
  value,
  className,
}: {
  value?: string | null;
  className?: string;
}) => {
  if (!value) return <span>—</span>;
  return (
    <code className={cn("text-muted-foreground", className)} title={value}>
      {value.slice(0, 8)}…
    </code>
  );
};
