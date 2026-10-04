import { useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { LucideCheck, LucideCopy } from "lucide-react";
import { cn } from "cn";

export interface CopyToClipboardButtonProps {
  value: string;
  size?: "default" | "lg";
}

export const CopyToClipboardButton = ({
  value,
  size = "default",
}: CopyToClipboardButtonProps) => {
  const feedbackDuration = 750;
  const [copied, setCopied] = useState(false);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const handleClick = () => {
    navigator.clipboard
      .writeText(value)
      .then(() => {
        setCopied(true);
        if (timeoutRef.current) clearTimeout(timeoutRef.current);
        timeoutRef.current = setTimeout(
          () => setCopied(false),
          feedbackDuration,
        );
      })
      .catch((err) => {
        console.error("Failed to copy text: ", err);
      });
  };

  const isLarge = size === "lg";

  return (
    <Button
      variant="ghost"
      size={isLarge ? "icon-lg" : "default"}
      onClick={handleClick}
      className="relative overflow-hidden"
    >
      <span
        className={cn(
          "absolute transition-all duration-200",
          copied ? "scale-0 opacity-0" : "scale-100 opacity-100",
        )}
      >
        <LucideCopy className={cn(isLarge && "size-5")} />
      </span>
      <span
        className={cn(
          "absolute transition-all duration-200",
          copied ? "scale-100 opacity-100" : "scale-0 opacity-0",
        )}
      >
        <LucideCheck className={cn("text-lm-status-up", isLarge && "size-5")} />
      </span>
    </Button>
  );
};
