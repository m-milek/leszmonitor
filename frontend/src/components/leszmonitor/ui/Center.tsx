import { cn } from "@/lib/utils.ts";

export interface CenterProps {
  children: React.ReactNode;
  centerVertically?: boolean;
  centerHorizontally?: boolean;
  className?: string;
}

export function Center({
  children,
  centerVertically = true,
  centerHorizontally = true,
  className,
}: Readonly<CenterProps>) {
  return (
    <div
      className={cn(
        `flex w-full h-full ${centerVertically ? "items-center" : ""} ${centerHorizontally ? "justify-center" : ""}`,
        className ?? "",
      )}
    >
      {children}
    </div>
  );
}
