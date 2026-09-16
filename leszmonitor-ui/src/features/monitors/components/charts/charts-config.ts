export const CHART_CONFIG = {
  VISIBLE_POINTS: 20,
  MARGIN: { left: -30, top: 5, right: 15 },
  X_AXIS: {
    minTickGap: 20,
  },
  ANIMATION: {
    duration: 800,
    easing: "linear" as const,
  },
} as const;
