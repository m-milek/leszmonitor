import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { Divider } from "@/components/leszmonitor/ui/Divider.tsx";

describe("Divider", () => {
  it("renders a horizontal divider when direction is 'row'", () => {
    render(<Divider direction="row" />);
    const divider = screen.getByRole("separator");

    expect(divider).toBeInTheDocument();
    expect(divider).toHaveClass("w-full", "border-t");
  });

  it("renders a vertical divider when direction is 'column'", () => {
    render(<Divider direction="column" />);
    const divider = screen.getByRole("separator");

    expect(divider).toBeInTheDocument();
    expect(divider).toHaveClass("border-l");
    expect(divider).not.toHaveClass("w-full");
  });

  it("applies additional class names passed via className prop", () => {
    render(<Divider direction="row" className="my-4" />);
    const divider = screen.getByRole("separator");

    expect(divider).toHaveClass("my-4", "border-t");
  });
});
