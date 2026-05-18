import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { UserInitial } from "./UserInitial";

describe("UserInitial", () => {
  it("renders the first letter uppercased", () => {
    render(<UserInitial username="admin" />);
    expect(screen.getByText("A")).toBeInTheDocument();
  });

  it("renders '?' for empty username", () => {
    render(<UserInitial username="" />);
    expect(screen.getByText("?")).toBeInTheDocument();
  });

  it("renders '?' for undefined-ish username", () => {
    render(<UserInitial username={undefined as unknown as string} />);
    expect(screen.getByText("?")).toBeInTheDocument();
  });

  it("applies the correct size class", () => {
    const { container } = render(<UserInitial username="bob" size="sm" />);
    expect(container.firstChild).toHaveClass("size-8");
  });

  it("defaults to xl size", () => {
    const { container } = render(<UserInitial username="bob" />);
    expect(container.firstChild).toHaveClass("size-24");
  });

  it("passes through additional className", () => {
    const { container } = render(
      <UserInitial username="bob" className="mt-4" />,
    );
    expect(container.firstChild).toHaveClass("mt-4");
  });
});
