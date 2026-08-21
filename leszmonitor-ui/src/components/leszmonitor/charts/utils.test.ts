import { describe, expect, it } from "vitest";
import { generateValuesWithInterval, padZero } from "./utils";

describe("generateValuesWithInterval", () => {
  it("returns empty array for length 0", () => {
    expect(generateValuesWithInterval(0, 5)).toEqual([]);
  });

  it("generates correct values", () => {
    expect(generateValuesWithInterval(4, 10)).toEqual([0, 10, 20, 30]);
  });

  it("handles interval of 1", () => {
    expect(generateValuesWithInterval(3, 1)).toEqual([0, 1, 2]);
  });
});

describe("padZero", () => {
  it("pads single digit with zero", () => {
    expect(padZero(0)).toBe("00");
    expect(padZero(5)).toBe("05");
    expect(padZero(9)).toBe("09");
  });

  it("does not pad double digits", () => {
    expect(padZero(10)).toBe("10");
    expect(padZero(99)).toBe("99");
  });
});
