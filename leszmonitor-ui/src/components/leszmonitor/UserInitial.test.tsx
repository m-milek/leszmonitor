import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { Initial } from "./Initial.tsx";

describe("UserInitial", () => {
  it("renders the first letter uppercased", () => {
    render(<Initial text="admin" />);
    expect(screen.getByText("A")).toBeInTheDocument();
  });

  it("renders '?' for empty username", () => {
    render(<Initial text="" />);
    expect(screen.getByText("?")).toBeInTheDocument();
  });

  it("applies the correct size class", () => {
    const { container } = render(<Initial text="bob" size="sm" />);
    expect(container.firstChild).toHaveClass("size-8");
  });

  it("defaults to xl size", () => {
    const { container } = render(<Initial text="bob" />);
    expect(container.firstChild).toHaveClass("size-24");
  });

  it("passes through additional className", () => {
    const { container } = render(<Initial text="bob" className="mt-4" />);
    expect(container.firstChild).toHaveClass("mt-4");
  });
});
