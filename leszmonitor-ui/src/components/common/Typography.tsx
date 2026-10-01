import type { ReactNode } from "react";

interface TypographyProps {
  readonly children: ReactNode;
}

export function TypographyH1({ children }: TypographyProps) {
  return <h1 className="scroll-m-20 text-3xl font-extrabold">{children}</h1>;
}

export function TypographyH2({ children }: TypographyProps) {
  return (
    <h2 className="scroll-m-20 text-xl font-semibold first:mt-0">{children}</h2>
  );
}

export function TypographyH3({ children }: TypographyProps) {
  return <h3 className="scroll-m-20 text-lg font-semibold">{children}</h3>;
}

export function TypographyH4({ children }: TypographyProps) {
  return <h4 className="scroll-m-20 font-semibold">{children}</h4>;
}
