import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { ErrorTooltip } from "./ErrorTooltip";
import { TooltipProvider } from "@/components/ui/tooltip";

function renderWithProvider(ui: React.ReactElement) {
  return render(<TooltipProvider>{ui}</TooltipProvider>);
}

describe("ErrorTooltip", () => {
  it("renders children", () => {
    renderWithProvider(
      <ErrorTooltip>
        <button>Click me</button>
      </ErrorTooltip>,
    );

    expect(screen.getByText("Click me")).toBeInTheDocument();
  });

  it("does NOT show tooltip content when isOpen is false", () => {
    renderWithProvider(
      <ErrorTooltip isOpen={false} message="Something went wrong">
        <button>Click me</button>
      </ErrorTooltip>,
    );

    expect(screen.queryByText("Something went wrong")).not.toBeInTheDocument();
  });

  it("shows tooltip content when isOpen is true", () => {
    renderWithProvider(
      <ErrorTooltip isOpen={true} message="Field is required">
        <button>Click me</button>
      </ErrorTooltip>,
    );

    // ✅ Use toHaveTextContent instead of toContain
    expect(screen.getByRole("tooltip")).toHaveTextContent("Field is required");
  });

  it("applies destructive styling to tooltip content", () => {
    renderWithProvider(
      <ErrorTooltip isOpen={true} message="Error!">
        <button>Click me</button>
      </ErrorTooltip>,
    );

    // ✅ Use selector to target the visible element, not the hidden a11y span
    const content = screen.getByText("Error!", {
      selector: "[data-slot='tooltip-content']",
    });
    expect(content).toHaveClass("bg-destructive");
  });

  it("renders with default props without crashing", () => {
    renderWithProvider(
      <ErrorTooltip>
        <span>Child</span>
      </ErrorTooltip>,
    );

    expect(screen.getByText("Child")).toBeInTheDocument();
  });

  it("wraps children in a span (for asChild)", () => {
    renderWithProvider(
      <ErrorTooltip>
        <button>Click me</button>
      </ErrorTooltip>,
    );

    const button = screen.getByText("Click me");
    expect(button.closest("span")).toBeInTheDocument();
  });
});
