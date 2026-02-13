import type { ReactNode } from "react";

interface TypographyProps {
  children: ReactNode;
}

export function TypographyH1({ children }: TypographyProps) {
  return (
    <h1 className="scroll-m-20 text-2xl font-extrabold tracking-tight text-balance dark:text-foreground">
      {children}
    </h1>
  );
}

export function TypographyH2({ children }: TypographyProps) {
  return (
    <h2 className="scroll-m-20 text-xl font-semibold tracking-tight first:mt-0">
      {children}
    </h2>
  );
}

export function TypographyH3({ children }: TypographyProps) {
  return (
    <h3 className="scroll-m-20 text-lg font-semibold tracking-tight">
      {children}
    </h3>
  );
}

export function TypographyH4({ children }: TypographyProps) {
  return (
    <h4 className="scroll-m-20 font-semibold tracking-tight">{children}</h4>
  );
}

export function TypographyP({ children }: TypographyProps) {
  return <p className="leading-7 not-first:mt-6">{children}</p>;
}
