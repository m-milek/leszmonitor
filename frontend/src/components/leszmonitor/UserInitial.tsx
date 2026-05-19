import { cn } from "@/lib/utils";

export interface UserInitialProps {
  username: string;
  size?: "sm" | "md" | "lg" | "xl";
  className?: string;
}

const sizeClasses = {
  sm: "size-8 text-[8px]",
  md: "size-12 text-[12px]",
  lg: "size-16 text-[16px]",
  xl: "size-24 text-[24px]",
};

export const UserInitial = ({
  username,
  size = "xl",
  className,
}: UserInitialProps) => {
  const value = username?.[0]?.toUpperCase() ?? "?";

  return (
    <div
      className={cn(
        "flex items-center justify-center rounded-full bg-primary select-none",
        sizeClasses[size],
        className,
      )}
    >
      <span className="text-[2em] leading-none text-white">{value}</span>
    </div>
  );
};
