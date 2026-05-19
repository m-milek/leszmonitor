import { render, screen, fireEvent, act } from "@testing-library/react";
import { describe, expect, it, vi, beforeEach } from "vitest";
import { CopyToClipboardButton } from "./CopyToClipboardButton";

describe("CopyToClipboardButton", () => {
  beforeEach(() => {
    Object.assign(navigator, {
      clipboard: {
        writeText: vi.fn().mockResolvedValue(undefined),
      },
    });
  });

  it("renders a button", () => {
    render(<CopyToClipboardButton value="hello" />);
    expect(screen.getByRole("button")).toBeInTheDocument();
  });

  it("copies value to clipboard on click", async () => {
    render(<CopyToClipboardButton value="hello" />);
    await act(async () => {
      fireEvent.click(screen.getByRole("button"));
    });
    expect(navigator.clipboard.writeText).toHaveBeenCalledWith("hello");
  });

  it("shows check icon after copying", async () => {
    const { container } = render(<CopyToClipboardButton value="test" />);
    const [copySpan, checkSpan] = container.querySelectorAll("span");

    expect(copySpan).toHaveClass("scale-100", "opacity-100");
    expect(checkSpan).toHaveClass("scale-0", "opacity-0");

    await act(async () => {
      fireEvent.click(screen.getByRole("button"));
    });

    expect(copySpan).toHaveClass("scale-0", "opacity-0");
    expect(checkSpan).toHaveClass("scale-100", "opacity-100");
  });

  it("handles clipboard failure gracefully", async () => {
    const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});
    vi.mocked(navigator.clipboard.writeText).mockRejectedValueOnce(
      new Error("denied"),
    );

    const { container } = render(<CopyToClipboardButton value="test" />);
    const [copySpan] = container.querySelectorAll("span");

    await act(async () => {
      fireEvent.click(screen.getByRole("button"));
    });

    expect(copySpan).toHaveClass("scale-100", "opacity-100");
    expect(consoleSpy).toHaveBeenCalled();
    consoleSpy.mockRestore();
  });
});
