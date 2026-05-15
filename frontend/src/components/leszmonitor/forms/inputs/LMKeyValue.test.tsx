import { render, screen } from "@testing-library/react";
import { describe, it, expect, vi } from "vitest";
import { TooltipProvider } from "@/components/ui/tooltip";
import {
  LMKeyValueInput,
  pairsToRecord,
} from "@/components/leszmonitor/forms/inputs/LMKeyValue.tsx";
import { userEvent } from "@testing-library/user-event/dist/cjs/setup/index.js";

function renderWithProvider(ui: React.ReactElement) {
  return render(<TooltipProvider>{ui}</TooltipProvider>);
}

describe("LMKeyValueInput", () => {
  it("renders initial key-value pairs from value prop", () => {
    renderWithProvider(
      <LMKeyValueInput
        name="headers"
        value={{ Authorization: "Bearer token", Accept: "application/json" }}
        onChange={() => {}}
      />,
    );

    expect(screen.getByDisplayValue("Authorization")).toBeInTheDocument();
    expect(screen.getByDisplayValue("Bearer token")).toBeInTheDocument();
    expect(screen.getByDisplayValue("Accept")).toBeInTheDocument();
    expect(screen.getByDisplayValue("application/json")).toBeInTheDocument();
  });

  it("renders empty state with just an Add button", () => {
    renderWithProvider(<LMKeyValueInput name="headers" onChange={() => {}} />);

    expect(screen.getByText("Add")).toBeInTheDocument();
    expect(screen.queryByPlaceholderText("Key")).not.toBeInTheDocument();
  });

  it("adds a new empty row when Add is clicked", async () => {
    const user = userEvent.setup();

    renderWithProvider(<LMKeyValueInput name="headers" onChange={() => {}} />);

    await user.click(screen.getByText("Add"));

    expect(screen.getByPlaceholderText("Key")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Value")).toBeInTheDocument();
  });

  it("calls onChange with updated record when a key is typed", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();

    renderWithProvider(<LMKeyValueInput name="headers" onChange={onChange} />);

    await user.click(screen.getByText("Add"));
    await user.type(screen.getByPlaceholderText("Key"), "Host");

    // onChange is called on each keystroke
    expect(onChange).toHaveBeenLastCalledWith({ Host: "" });
  });

  it("calls onChange with full record when key and value are typed", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();

    renderWithProvider(<LMKeyValueInput name="headers" onChange={onChange} />);

    await user.click(screen.getByText("Add"));
    await user.type(screen.getByPlaceholderText("Key"), "Host");
    await user.type(screen.getByPlaceholderText("Value"), "example.com");

    expect(onChange).toHaveBeenLastCalledWith({ Host: "example.com" });
  });

  it("removes a row when delete button is clicked", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();

    renderWithProvider(
      <LMKeyValueInput
        name="headers"
        value={{ Host: "example.com" }}
        onChange={onChange}
      />,
    );

    expect(screen.getByDisplayValue("Host")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "Remove" }));

    expect(screen.queryByDisplayValue("Host")).not.toBeInTheDocument();
    expect(onChange).toHaveBeenLastCalledWith({});
  });

  it("filters out pairs with empty key and value from onChange output", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();

    renderWithProvider(
      <LMKeyValueInput
        name="headers"
        value={{ Host: "example.com" }}
        onChange={onChange}
      />,
    );

    // Add an empty row - shouldn't appear in output
    await user.click(screen.getByText("Add"));

    // Edit the existing row to trigger onChange
    await user.type(screen.getByDisplayValue("example.com"), "!");

    expect(onChange).toHaveBeenLastCalledWith({ Host: "example.com!" });
  });

  it("uses custom placeholders and button text", () => {
    renderWithProvider(
      <LMKeyValueInput
        name="env"
        value={{ NODE_ENV: "prod" }}
        onChange={() => {}}
        keyPlaceholder="Variable"
        valuePlaceholder="Content"
        addButtonText="Add variable"
      />,
    );

    expect(screen.getByPlaceholderText("Variable")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Content")).toBeInTheDocument();
    expect(screen.getByText("Add variable")).toBeInTheDocument();
  });
});

describe("pairsToRecord", () => {
  it("filters out empty pairs", () => {
    expect(
      pairsToRecord([
        { id: "1", key: "", value: "" },
        { id: "2", key: "a", value: "b" },
      ]),
    ).toEqual({ a: "b" });
  });
});
