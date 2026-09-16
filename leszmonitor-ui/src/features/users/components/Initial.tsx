import { cn } from "@/lib/utils";
import { colorFromString } from "@/lib/colorFromString.ts";

export interface UserInitialProps {
  text: string;
  textForColorCalculation?: string;
  size?: "sm" | "md" | "lg" | "xl";
  className?: string;
}

const sizeClasses = {
  sm: "size-8 text-[8px]",
  md: "size-12 text-[12px]",
  lg: "size-16 text-[16px]",
  xl: "size-24 text-[24px]",
};

export const Initial = ({
  text,
  textForColorCalculation,
  size = "xl",
  className,
}: UserInitialProps) => {
  const value = text?.[0]?.toUpperCase() ?? "?";

  const backgroundColor = colorFromString(textForColorCalculation ?? text);

  return (
    <div
      className={cn(
        "flex items-center justify-center rounded-full bg-primary select-none",
        sizeClasses[size],
        className,
      )}
      style={{ backgroundColor }}
    >
      <span className="text-[2em] leading-none text-slate-800">{value}</span>
    </div>
  );
};
