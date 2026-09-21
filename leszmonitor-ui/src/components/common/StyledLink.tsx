import { Link, type LinkComponentProps } from "@tanstack/react-router";
import { cn } from "cn";
import * as React from "react";

export interface LinkProps extends LinkComponentProps {
  children: React.ReactNode;
}

export const StyledLink = (props: LinkProps) => {
  return (
    <Link
      {...props}
      className={cn("text-primary hover:underline", props.className)}
    >
      {props.children}
    </Link>
  );
};
