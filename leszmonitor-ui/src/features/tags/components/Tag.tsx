import type { Tag as TagModel } from "@/features/tags/model/types.ts";
import { Badge } from "@/components/ui/badge.tsx";
import { cn } from "@/lib/utils.ts";
import { tagChipStyle } from "@/features/tags/lib/colors.ts";

export interface TagProps {
  tag: Pick<TagModel, "name" | "colorHex"> &
    Partial<Pick<TagModel, "description">>;
  className?: string;
}

export const Tag = ({ tag, className }: TagProps) => (
  <Badge
    title={tag.description}
    style={tagChipStyle(tag.colorHex)}
    className={cn("py-2 px-3 max-w-full truncate border", className)}
  >
    {tag.name}
  </Badge>
);
