import { describe, it, expect } from "vitest";
import {
  hslToHex,
  isValidHexColor,
  normalizeHexColor,
  randomTagColor,
  tagChipStyle,
} from "@/lib/tagColors.ts";

describe("hslToHex", () => {
  const tests: [[number, number, number], string][] = [
    [[0, 100, 50], "#ff0000"],
    [[120, 100, 50], "#00ff00"],
    [[240, 100, 50], "#0000ff"],
    [[0, 0, 100], "#ffffff"],
    [[0, 0, 0], "#000000"],
    [[0, 0, 50], "#808080"],
  ];

  it.each(tests)("converts hsl%s to %s", ([h, s, l], expected) => {
    expect(hslToHex(h, s, l)).toBe(expected);
  });
});

describe("isValidHexColor", () => {
  it.each(["#fff", "#FFF", "#3b82f6", "#3B82F6"])("accepts %s", (color) => {
    expect(isValidHexColor(color)).toBe(true);
  });

  it.each(["fff", "#ffff", "#gggggg", "", "rgb(0,0,0)"])(
    "rejects %s",
    (color) => {
      expect(isValidHexColor(color)).toBe(false);
    },
  );
});

describe("normalizeHexColor", () => {
  const tests: [string, string][] = [
    ["#ABC", "#aabbcc"],
    ["#abc", "#aabbcc"],
    ["#3B82F6", "#3b82f6"],
    ["  #3b82f6  ", "#3b82f6"],
  ];

  it.each(tests)("normalizes %s to %s", (input, expected) => {
    expect(normalizeHexColor(input)).toBe(expected);
  });
});

describe("randomTagColor", () => {
  it("generates valid hex colors", () => {
    for (let i = 0; i < 100; i++) {
      expect(isValidHexColor(randomTagColor())).toBe(true);
    }
  });
});

describe("tagChipStyle", () => {
  it("uses the same accent for the border and the label", () => {
    const style = tagChipStyle("#3b82f6");

    expect(style.color).toBe(style.borderColor);
    expect(style.color).toContain("#3b82f6");
    expect(style.backgroundColor).toContain("#3b82f6");
    expect(style.backgroundColor).not.toBe(style.color);
  });

  it("normalizes the shorthand form", () => {
    expect(tagChipStyle("#ABC").color).toContain("#aabbcc");
  });

  it("falls back to a neutral color for invalid input", () => {
    expect(tagChipStyle("nope").color).toContain("#6b7280");
  });
});
