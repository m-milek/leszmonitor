import { RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button.tsx";
import { LMInputField } from "@/components/form/LMInputField.tsx";
import { Flex } from "@/components/common/Flex.tsx";
import { randomTagColor, tagChipStyle } from "@/features/tags/lib/colors.ts";

export interface LMColorPickerProps {
  name: string;
  value: string;
  onChange: (colorHex: string) => void;
  isInvalid?: boolean;
  errorMessage?: string;
}

export const LMColorPicker = ({
  name,
  value,
  onChange,
  isInvalid,
  errorMessage,
}: LMColorPickerProps) => (
  <Flex className="items-center gap-2">
    <Button
      type="button"
      variant="outline"
      size="icon"
      title="Pick a random color"
      aria-label="Pick a random color"
      onClick={() => onChange(randomTagColor())}
      style={tagChipStyle(value)}
      className="shrink-0 border"
    >
      <RefreshCw className="h-4 w-4" />
    </Button>
    <LMInputField
      name={name}
      type="text"
      value={value}
      onChange={(e) => onChange(e.target.value)}
      isInvalid={isInvalid}
      errorMessage={errorMessage}
    />
  </Flex>
);
