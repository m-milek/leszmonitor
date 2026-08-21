export const colorFromString = (input: string): string => {
  let hash = 0;
  for (let i = 0; i < input.length; i++) {
    hash = input.charCodeAt(i) + ((hash << 5) - hash);
  }
  const hue = Math.abs(hash) % 360; // Return a hue value between 0 and 359

  return `hsl(${hue}, 75%, 70%)`;
};
