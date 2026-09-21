import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { MonitorStatusPill } from "./MonitorStatusPill";
import type { Monitor } from "@/features/monitors/types";

const createMonitor = (state: string): Monitor =>
  ({
    runState: state,
  }) as Monitor;

describe("MonitorStatusPill", () => {
  it("renders 'Active' with green styling when monitor state is active", () => {
    render(<MonitorStatusPill monitor={createMonitor("active")} />);

    const pill = screen.getByText("Active");
    expect(pill).toBeInTheDocument();
  });

  it("renders 'Paused' with gray styling when monitor state is paused", () => {
    render(<MonitorStatusPill monitor={createMonitor("paused")} />);

    const pill = screen.getByText("Paused");
    expect(pill).toBeInTheDocument();
  });

  it("renders 'Invalid' with muted styling for an unknown monitor state", () => {
    render(<MonitorStatusPill monitor={createMonitor("unknown")} />);

    const pill = screen.getByText("Invalid");
    expect(pill).toBeInTheDocument();
  });

  it("renders a Badge coloured from the status tokens", () => {
    render(<MonitorStatusPill monitor={createMonitor("active")} />);

    const pill = screen.getByText("Active");
    expect(pill).toHaveAttribute("data-slot", "badge");
    expect(pill.className).toContain("bg-lm-status-up");
  });
});
