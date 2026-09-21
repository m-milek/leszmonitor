import { cn } from "cn";
import { avatarStyleFromString } from "@/lib/avatar-color";

export interface UserInitialProps {
  text: string;
  textForColorCalculation?: string;
  size?: "sm" | "md" | "lg" | "xl";
  className?: string;
}

const sizeClasses = {
  sm: "size-8 text-base",
  md: "size-12 text-2xl",
  lg: "size-16 text-3xl",
  xl: "size-24 text-5xl",
};

export const Initial = ({
  text,
  textForColorCalculation,
  size = "xl",
  className,
}: UserInitialProps) => {
  const value = text?.[0]?.toUpperCase() ?? "?";

  return (
    <div
      className={cn(
        "flex items-center justify-center rounded-full font-medium select-none",
        sizeClasses[size],
        className,
      )}
      style={avatarStyleFromString(textForColorCalculation ?? text)}
    >
      <span className="leading-none">{value}</span>
    </div>
  );
};
