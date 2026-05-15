import { describe, it, expect } from "vitest";
import { slugFromString } from "@/lib/slugFromString.ts";

describe("slugFromString", () => {
  const tests: [string, string][] = [
    ["Hello World", "hello-world"],
    ["already-hyphenated", "already-hyphenated"],
    ["Mix-ed Hyphen", "mix-ed-hyphen"],
    ["Special@Chars#Here!", "specialcharshere"],
    ["  Trim  Spaces  ", "trim-spaces"],
    ["Multiple   Spaces", "multiple-spaces"],
    ["UPPERCASE", "uppercase"],
    ["123-numbers-456", "123-numbers-456"],
    ["", ""],
    ["---", ""],
    ["Test - With - Dashes", "test-with-dashes"],
    ["consecutive---hyphens", "consecutive-hyphens"],
    ["@#$%^&*()", ""],
    ["   ", ""],
    ["-leading-hyphen", "leading-hyphen"],
    ["trailing-hyphen-", "trailing-hyphen"],
    ["one-two--three---four", "one-two-three-four"],
    ["CamelCase", "camelcase"],
    ["snake_case", "snakecase"],
    ["dot.separated.words", "dotseparatedwords"],
    ["email@example.com", "emailexamplecom"],
    ["tabs\there\ttoo", "tabs-here-too"],
    ["newline\ntest", "newline-test"],
    ["!start-middle-end!", "start-middle-end"],
    ["a", "a"],
    ["1", "1"],
    ["-", ""],
    ["test\u2014with\u2014em\u2014dash", "testwithemdash"],
    ["test\u2013with\u2013en\u2013dash", "testwithendash"],
  ];

  for (const [input, expected] of tests) {
    it(`"${input}" -> "${expected}"`, () => {
      expect(slugFromString(input)).toBe(expected);
    });
  }
});
