# collapsible

2026-09-20, transformation engine (legacy `new-york` style: classified against
the stock `new-york-v4` golden, then transformed in place). Migrated to
`@base-ui/react/collapsible`; `Content` → `Panel`.

## Changed

- `src/components/ui/collapsible.tsx` — `import { Collapsible as CollapsiblePrimitive } from "radix-ui"`
  → `from "@base-ui/react/collapsible"`. Part rewiring:
  `CollapsiblePrimitive.CollapsibleTrigger` → `CollapsiblePrimitive.Trigger`,
  `CollapsiblePrimitive.CollapsibleContent` → `CollapsiblePrimitive.Panel`;
  `Root` unchanged. Props types →
  `CollapsiblePrimitive.Root.Props` / `.Trigger.Props` / `.Panel.Props`
  (the file previously referenced `React.ComponentProps` without importing
  React, relying on the global namespace; that dependency is gone).
  The three public export names (`Collapsible`, `CollapsibleTrigger`,
  `CollapsibleContent`) and all three `data-slot` values are deliberately kept,
  including `data-slot="collapsible-content"` on the Panel, so any CSS hook or
  test selector keeps working.
  Leftover scan: `grep -n "radix-ui\|@radix-ui" src/components/ui/collapsible.tsx` → clean.

Classification: PRISTINE against `new-york-v4`; the result matches the
`base-nova` registry collapsible structurally.

## Left alone

Nothing to sweep: this wrapper has no consumers anywhere in `src` (verified by
`grep -rn 'ui/collapsible' src`). It is an unused shadcn install, migrated so the
project can drop Radix, but no call site exercises it.

## Behavior changes

- `forceMount` is now `keepMounted` on the Panel, and `hiddenUntilFound` is
  newly available. Nothing in the project passes either.
- The Trigger's open marker moved from `data-state="open"` to `data-panel-open`,
  and Panel state moved to `data-open` / `data-closed`. This wrapper ships no
  classes, so there was nothing to rewrite — but any future consumer styling the
  trigger must use `data-panel-open:`, not `data-open:`.
- Collapsible parts no longer emit `data-disabled`; gate on the `disabled` prop
  or `:disabled` on the trigger instead.
- Panel height animations use `--collapsible-panel-height` (was
  `--radix-collapsible-content-height`) with `data-starting-style` /
  `data-ending-style` transitions rather than `animate-in` / `animate-out`.

## Verify by hand

Nothing renders this component today, so there is no in-app path to click. If
you add a consumer, check: trigger toggles the panel, Space/Enter works,
`data-panel-open` drives any chevron rotation, and the panel does not flash at
full height on first paint.
