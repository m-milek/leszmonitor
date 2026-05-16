export const generateValuesWithInterval = (
  length: number,
  interval: number,
): number[] => {
  return Array.from({ length }, (_, i) => i * interval);
};

const padZero = (x: number): string => {
  return x < 10 ? `0${x}` : `${x}`;
};

export const formatTime = (timestamp: string): string => {
  const date = new Date(timestamp);
  return `${padZero(date.getHours())}:${padZero(date.getMinutes())}:${padZero(date.getSeconds())}`;
};
