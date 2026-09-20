# separator

2026-09-20, transformation engine (legacy `new-york` style: classified against
the stock `new-york-v4` golden, then transformed in place). Migrated to the
callable `@base-ui/react/separator` primitive.

## Changed

- `src/components/ui/separator.tsx` — `import { Separator as SeparatorPrimitive } from "radix-ui"`
  → `from "@base-ui/react/separator"`; `<SeparatorPrimitive.Root>` →
  callable `<SeparatorPrimitive>` (single-part primitive, no `.Root`).
  Props type `React.ComponentProps<typeof SeparatorPrimitive.Root>` →
  `SeparatorPrimitive.Props`; the now-unused `import * as React` was removed.
  The `decorative` prop (and its `decorative={decorative}` forward) is gone —
  dropped by Base UI, whose separator is always semantic `role="separator"`.
  `data-slot="separator"`, `orientation` defaulting to `"horizontal"`, the
  `"use client"` directive and the class string are unchanged;
  `data-[orientation=...]` needs no rewrite (still parameterised in Base UI).
  Leftover scan: `grep -n "radix-ui\|@radix-ui" src/components/ui/separator.tsx` → clean.

Classification: PRISTINE against `new-york-v4` (class-order only, from the
project's Prettier Tailwind ordering).

## Left alone

- `src/components/ui/button-group.tsx:66`, `src/components/ui/sidebar.tsx:362`,
  `src/components/ui/field.tsx:179` — the three consumers. None passed
  `decorative`, so no call site changed. They pass only `orientation`,
  `className` and `data-*`, all still valid.

## Behavior changes

- The separator is now always exposed to assistive tech as `role="separator"`.
  Under Radix these three call sites were implicitly `decorative={true}`, i.e.
  `role="none"` and `aria-hidden`. Screen readers will now announce a separator
  at each of them. Flagged, not patched: Base UI has no `decorative` equivalent;
  a purely visual rule would have to become a plain `<div aria-hidden="true">`,
  which is a design decision, not a migration step.

## Verify by hand

1. Sidebar: the divider under the header still spans the full width with the
   `mx-2` inset and `bg-sidebar-border` colour.
2. A `ButtonGroup` with `orientation="vertical"`: the separator must stretch the
   full group height (`data-[orientation=vertical]:h-full`), not collapse to 0.
3. `FieldSeparator` (field.tsx): the rule still sits centred behind its label
   text via `absolute inset-0 top-1/2`.
