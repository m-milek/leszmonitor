import type { Tag as TagModel } from "@/features/tags/types";
import { cn } from "cn";
import { Badge } from "@/components/ui/badge";
import {
  BADGE_SIZE_CLASS,
  type BadgeSize,
} from "@/components/common/badge-size";
import { tagChipStyle } from "@/features/tags/lib/colors";

type TagDisplay = Pick<TagModel, "name" | "colorHex"> &
  Partial<Pick<TagModel, "description">>;

export interface TagProps {
  tag: TagDisplay;
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
