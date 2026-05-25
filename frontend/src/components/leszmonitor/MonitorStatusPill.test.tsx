import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { MonitorStatusPill } from "./MonitorStatusPill";
import type { Monitor } from "@/lib/types.ts";

const createMonitor = (state: string): Monitor =>
  ({
    state,
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

  it("renders as a span with rounded pill classes", () => {
    render(<MonitorStatusPill monitor={createMonitor("active")} />);

    const pill = screen.getByText("Active");
    expect(pill.tagName).toBe("SPAN");
    expect(pill.className).toContain("inline-flex");
    expect(pill.className).toContain("rounded-xl");
    expect(pill.className).toContain("font-medium");
  });
});
