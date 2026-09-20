# label

2026-09-20, transformation engine (legacy `new-york` style: classified against
the stock `new-york-v4` golden, then transformed in place). Migrated; Label has
no Base UI counterpart, so it is now a native `<label>`.

## Changed

- `src/components/ui/label.tsx` — dropped `import { Label as LabelPrimitive } from "radix-ui"`;
  `<LabelPrimitive.Root>` → native `<label>` (`display-misc.md`: no Base UI
  counterpart; `Field.Label` only applies inside a `Field.Root`). Props type
  `React.ComponentProps<typeof LabelPrimitive.Root>` → `React.ComponentProps<"label">`.
  `data-slot="label"` and the class string are unchanged — `select-none` already
  carried Radix's only behavioral extra (no text selection on double click).
  Leftover scan: `grep -n "radix-ui\|@radix-ui" src/components/ui/label.tsx` → clean.

Classification: PRISTINE against `new-york-v4` (only the project's `cn` alias
and the removed `"use client"` differ), and the resulting file is class-for-class
identical to the `base-nova` registry label.

## Left alone

- `src/components/ui/field.tsx:120` — the only consumer. Passes `className`,
  `data-slot` and `htmlFor`, all valid on a native `<label>`; no call-site change
  needed.

## Behavior changes

None. `htmlFor`, `onClick` and the peer/group class hooks behave identically; the
element rendered was already `<label>` under Radix.

## Verify by hand

1. Open any form (Monitors → New Monitor). Click a field's label text — focus
   should jump into the matching input.
2. Double-click a label — the text must not become selected.
3. Disable a field and confirm the label still dims (`peer-disabled:opacity-50`).
