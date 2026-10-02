import type { Tag as TagModel } from "@/features/tags/types";
import { cn } from "cn";
import { Badge } from "@/components/ui/badge";
import {
  BADGE_SIZE_CLASS,
  type BadgeSize,
} from "@/components/common/badge-size";
import { tagChipStyle } from "@/features/tags/lib/colors";

export interface TagProps {
  tag: TagModel;
  size?: BadgeSize;
}

export const Tag = ({ tag, size = "default" }: TagProps) => {
  return (
    <Badge
      title={tag.description}
      style={tagChipStyle(tag.colorHex)}
      className={cn(BADGE_SIZE_CLASS[size], size === "lg" && "w-18")}
    >
      {tag.name}
    </Badge>
  );
};
