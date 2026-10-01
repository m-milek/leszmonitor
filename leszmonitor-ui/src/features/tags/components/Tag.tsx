import type { Tag as TagModel } from "@/features/tags/types";
import { Badge } from "@/components/ui/badge";
import { tagChipStyle } from "@/features/tags/lib/colors";

export interface TagProps {
  tag: TagModel;
  size?: "default" | "lg";
}

export const Tag = ({ tag, size = "default" }: TagProps) => {
  const className = size === "default" ? "" : "text-sm h-8 w-18";
  return (
    <Badge
      title={tag.description}
      style={tagChipStyle(tag.colorHex)}
      className={className}
    >
      {tag.name}
    </Badge>
  );
};
