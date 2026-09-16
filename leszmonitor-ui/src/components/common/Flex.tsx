import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cn } from "@/lib/utils";

interface FlexProps extends React.HTMLAttributes<HTMLDivElement> {
  direction?: "row" | "column" | "row-reverse" | "column-reverse";
  asChild?: boolean;
}

const directionClass: Record<string, string> = {
  row: "flex-row",
  column: "flex-col",
  "row-reverse": "flex-row-reverse",
  "column-reverse": "flex-col-reverse",
};

const Flex = React.forwardRef<HTMLDivElement, FlexProps>(
  ({ direction = "row", asChild, className, ...props }, ref) => {
    const Comp = asChild ? Slot : "div";

    return (
      <Comp
        ref={ref}
        className={cn("flex", directionClass[direction], className)}
        {...props}
      />
    );
  },
);
Flex.displayName = "Flex";

export { Flex, type FlexProps };
