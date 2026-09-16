import { Link, type LinkComponentProps } from "@tanstack/react-router";
import * as React from "react";

export interface LinkProps extends LinkComponentProps {
  children: React.ReactNode;
}

export const StyledLink = (props: LinkProps) => {
  return (
    <Link
      {...props}
      className={`text-primary hover:underline ${props.className ?? ""}`}
    >
      {props.children}
    </Link>
  );
};
