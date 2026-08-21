import { useRef, useState } from "react";
import { Button } from "../ui/button";
import { LucideCheck, LucideCopy } from "lucide-react";

export interface CopyToClipboardButtonProps {
  value: string;
}

export const CopyToClipboardButton = ({
  value,
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

  return (
    <Button
      variant="ghost"
      onClick={handleClick}
      className="relative overflow-hidden"
    >
      <span
        className={`absolute transition-all duration-200 ${
          copied ? "scale-0 opacity-0" : "scale-100 opacity-100"
        }`}
      >
        <LucideCopy className="h-4 w-4" />
      </span>
      <span
        className={`absolute transition-all duration-200 ${
          copied ? "scale-100 opacity-100" : "scale-0 opacity-0"
        }`}
      >
        <LucideCheck className="h-4 w-4 text-green-500" />
      </span>
    </Button>
  );
};
