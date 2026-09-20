import { mergeProps } from "@base-ui/react/merge-props";
import { useRender } from "@base-ui/react/use-render";
import { cn } from "cn";

interface FlexProps extends useRender.ComponentProps<"div"> {
  direction?: "row" | "column" | "row-reverse" | "column-reverse";
}

const directionClass: Record<string, string> = {
  row: "flex-row",
  column: "flex-col",
  "row-reverse": "flex-row-reverse",
  "column-reverse": "flex-col-reverse",
};

function Flex({ direction = "row", className, render, ...props }: FlexProps) {
  return useRender({
    defaultTagName: "div",
    render,
    props: mergeProps<"div">(
      { className: cn("flex", directionClass[direction], className) },
      props,
    ),
  });
}

export { Flex, type FlexProps };
