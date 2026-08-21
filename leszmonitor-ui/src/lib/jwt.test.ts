import { describe, it, expect, vi } from "vitest";
import { isJwtClaims, isJwtValid } from "./jwt"; // adjust path

// fake JWT with a given payload
function createFakeJwt(payload: object): string {
  const header = btoa(JSON.stringify({ alg: "HS256" }));
  const body = btoa(JSON.stringify(payload));
  return `${header}.${body}.fake-signature`;
}

describe("isJwtClaims", () => {
  it("returns true for valid claims", () => {
    expect(isJwtClaims({ username: "alice", iat: 1000, exp: 2000 })).toBe(true);
  });

  it("returns false for null", () => {
    expect(isJwtClaims(null)).toBe(false);
  });

  it("returns false for non-object", () => {
    expect(isJwtClaims("string")).toBe(false);
    expect(isJwtClaims(42)).toBe(false);
    expect(isJwtClaims(undefined)).toBe(false);
  });

  it("returns false when username is missing", () => {
    expect(isJwtClaims({ iat: 1000, exp: 2000 })).toBe(false);
  });

  it("returns false when username is not a string", () => {
    expect(isJwtClaims({ username: 123, iat: 1000, exp: 2000 })).toBe(false);
  });

  it("returns false when exp is missing", () => {
    expect(isJwtClaims({ username: "alice", iat: 1000 })).toBe(false);
  });

  it("returns false when iat is not a number", () => {
    expect(isJwtClaims({ username: "alice", iat: "now", exp: 2000 })).toBe(
      false,
    );
  });
});

describe("isJwtValid", () => {
  it("returns claims for a valid, non-expired token", () => {
    const futureExp = Math.floor(Date.now() / 1000) + 3600; // 1 hour from now
    const token = createFakeJwt({
      username: "alice",
      iat: 1000,
      exp: futureExp,
    });

    const result = isJwtValid(token);

    expect(result).toEqual({ username: "alice", iat: 1000, exp: futureExp });
  });

  it("returns null for an expired token", () => {
    const pastExp = Math.floor(Date.now() / 1000) - 3600; // 1 hour ago
    const token = createFakeJwt({ username: "alice", iat: 1000, exp: pastExp });

    expect(isJwtValid(token)).toBeNull();
  });

  it("returns null when payload is missing required fields", () => {
    const token = createFakeJwt({ foo: "bar" });

    expect(isJwtValid(token)).toBeNull();
  });

  it("returns null for malformed base64", () => {
    expect(isJwtValid("not.valid-base64!.token")).toBeNull();
  });

  it("returns null for completely invalid input", () => {
    expect(isJwtValid("garbage")).toBeNull();
  });

  it("returns null for empty string", () => {
    expect(isJwtValid("")).toBeNull();
  });

  it("returns null when token has valid base64 but invalid JSON", () => {
    const header = btoa("not json");
    const body = btoa("also not json");
    expect(isJwtValid(`${header}.${body}.sig`)).toBeNull();
  });

  it("checks expiry against current time", () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2026-01-01T00:00:00Z"));

    const token = createFakeJwt({
      username: "alice",
      iat: 1000,
      exp: Math.floor(new Date("2026-06-01").getTime() / 1000),
    });

    expect(isJwtValid(token)).not.toBeNull();

    // Fast-forward past expiry
    vi.setSystemTime(new Date("2027-01-01T00:00:00Z"));
    expect(isJwtValid(token)).toBeNull();

    vi.useRealTimers();
  });
});
