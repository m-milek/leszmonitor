const HEX_COLOR_REGEX = /^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$/;

const RANDOM_SATURATION = { min: 58, max: 78 };
const RANDOM_LIGHTNESS = { min: 44, max: 62 };

const FALLBACK_COLOR = "#6b7280";
const BACKGROUND_MIX = "14%";
const ACCENT_MIX = "55%";

export const isValidHexColor = (value: string): boolean =>
  HEX_COLOR_REGEX.test(value);

export const normalizeHexColor = (value: string): string => {
  const color = value.trim().toLowerCase();

  if (color.length === 4) {
    return `#${color[1]}${color[1]}${color[2]}${color[2]}${color[3]}${color[3]}`;
  }

  return color;
};

export const hslToHex = (
  hue: number,
  saturation: number,
  lightness: number,
): string => {
  const s = saturation / 100;
  const l = lightness / 100;
  const chroma = (1 - Math.abs(2 * l - 1)) * s;

  const channel = (n: number) => {
    const k = (n + hue / 30) % 12;
    const value = l - (chroma / 2) * Math.max(-1, Math.min(k - 3, 9 - k, 1));
    return Math.round(value * 255)
      .toString(16)
      .padStart(2, "0");
  };

  return `#${channel(0)}${channel(8)}${channel(4)}`;
};

const randomBetween = (min: number, max: number) =>
  min + Math.random() * (max - min);

export const randomTagColor = (): string =>
  hslToHex(
    Math.random() * 360,
    randomBetween(RANDOM_SATURATION.min, RANDOM_SATURATION.max),
    randomBetween(RANDOM_LIGHTNESS.min, RANDOM_LIGHTNESS.max),
  );

export const tagChipStyle = (colorHex: string) => {
  const color = isValidHexColor(colorHex)
    ? normalizeHexColor(colorHex)
    : FALLBACK_COLOR;
  const accent = `color-mix(in srgb, ${color} ${ACCENT_MIX}, var(--foreground))`;

  return {
    backgroundColor: `color-mix(in srgb, ${color} ${BACKGROUND_MIX}, var(--background))`,
    borderColor: accent,
    color: accent,
  };
};
