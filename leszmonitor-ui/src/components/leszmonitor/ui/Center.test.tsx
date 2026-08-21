import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { Center } from "./Center";

describe("Center", () => {
  it("renders children inside a centered container", () => {
    render(
      <Center>
        <span>Centered content</span>
      </Center>,
    );

    const content = screen.getByText("Centered content");
    const container = content.parentElement;

    if (!container) {
      throw new Error("Expected container element.");
    }

    expect(container.className).toContain("items-center");
    expect(container.className).toContain("justify-center");
  });
});
