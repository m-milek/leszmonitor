import { mergeProps } from "@base-ui/react/merge-props";
import { useRender } from "@base-ui/react/use-render";
import { cn } from "cn";

interface FlexProps extends useRender.ComponentProps<"div"> {
  direction?: FlexDirection;
  directionMobile?: FlexDirection;
}

type FlexDirection = "row" | "column" | "row-reverse" | "column-reverse";

const directionClass: Record<FlexDirection, string> = {
  row: "flex-row",
  column: "flex-col",
  "row-reverse": "flex-row-reverse",
  "column-reverse": "flex-col-reverse",
};

const mobileDirectionClass: Record<FlexDirection, string> = {
  row: "max-md:flex-row",
  column: "max-md:flex-col",
  "row-reverse": "max-md:flex-row-reverse",
  "column-reverse": "max-md:flex-col-reverse",
};

function Flex({
  direction = "row",
  directionMobile = direction,
  className,
  render,
  ...props
}: FlexProps) {
  return useRender({
    defaultTagName: "div",
    render,
    props: mergeProps<"div">(
      {
        className: cn(
          "flex",
          directionClass[direction],
          directionMobile !== direction &&
            mobileDirectionClass[directionMobile],
          className,
        ),
      },
      props,
    ),
  });
}

export { Flex, type FlexProps };
