export const avatarStyleFromString = (
  input: string,
): { backgroundColor: string; color: string } => {
  let hash = 0;
  for (let i = 0; i < input.length; i++) {
    hash = input.charCodeAt(i) + ((hash << 5) - hash);
  }
  const hue = Math.abs(hash) % 360;

  return {
    backgroundColor: `hsl(${hue}, 75%, 70%)`,
    color: `hsl(${hue}, 75%, 20%)`,
  };
};
