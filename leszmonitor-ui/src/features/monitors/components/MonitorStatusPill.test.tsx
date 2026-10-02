import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { MonitorStatePill } from "./MonitorStatePill.tsx";
import type { Monitor } from "@/features/monitors/types";

const createMonitor = (state: string): Monitor =>
  ({
    runState: state,
  }) as Monitor;

describe("MonitorStatusPill", () => {
  it("renders 'Active' with green styling when monitor state is active", () => {
    render(<MonitorStatePill monitor={createMonitor("active")} />);

    const pill = screen.getByText("Active");
    expect(pill).toBeInTheDocument();
  });

  it("renders 'Paused' with gray styling when monitor state is paused", () => {
    render(<MonitorStatePill monitor={createMonitor("paused")} />);

    const pill = screen.getByText("Paused");
    expect(pill).toBeInTheDocument();
  });

  it("renders 'Invalid' with muted styling for an unknown monitor state", () => {
    render(<MonitorStatePill monitor={createMonitor("unknown")} />);

    const pill = screen.getByText("Invalid");
    expect(pill).toBeInTheDocument();
  });

  it("renders a Badge coloured from the status tokens", () => {
    render(<MonitorStatePill monitor={createMonitor("active")} />);

    const pill = screen.getByText("Active");
    expect(pill).toHaveAttribute("data-slot", "badge");
    expect(pill.className).toContain("bg-lm-status-up");
  });
});
