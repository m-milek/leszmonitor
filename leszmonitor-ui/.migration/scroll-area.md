# scroll-area

2026-09-20, transformation engine (legacy `new-york` style: classified against
the stock `new-york-v4` golden, then transformed in place — the file is
CUSTOMIZED, so the customization was preserved rather than overwritten).
Migrated to `@base-ui/react/scroll-area`.

## Changed

- `src/components/ui/scroll-area.tsx` — `import { ScrollArea as ScrollAreaPrimitive } from "radix-ui"`
  → `from "@base-ui/react/scroll-area"`. Part renames:
  `ScrollAreaPrimitive.ScrollAreaScrollbar` → `ScrollAreaPrimitive.Scrollbar`
  (`:36`, `:53`), `ScrollAreaPrimitive.ScrollAreaThumb` →
  `ScrollAreaPrimitive.Thumb` (`:49`); `Root`, `Viewport` and `Corner` keep their
  names. Props types → `ScrollAreaPrimitive.Root.Props` and
  `ScrollAreaPrimitive.Scrollbar.Props`; the now-unused `import * as React` was
  removed. No class string changed: this wrapper carried no `data-[state=…]`
  hooks, and `orientation` styling is done in JS
  (`orientation === "vertical" && …`), which is attribute-independent and
  therefore survives Base UI dropping Radix's `data-state` on the scrollbar.
  All five `data-slot` values are unchanged.
  PRESERVED CUSTOMIZATION (`:23`-`:24`): the project renders
  `<ScrollBar orientation="vertical" />` and `<ScrollBar orientation="horizontal" />`
  explicitly where both the `new-york-v4` and `base-nova` registry files render a
  single default `<ScrollBar />`. Both calls were kept verbatim.
  Leftover scan: `grep -n "radix-ui\|@radix-ui\|data-\[state=" src/components/ui/scroll-area.tsx` → clean.

`ScrollArea.Content` (a Base UI-only part that wraps children inside the
Viewport for horizontal overflow measurement) was deliberately NOT added: the
`base-nova` registry wrapper does not use it, and adding it would be a
structural change beyond the migration.

## Left alone

- `src/features/monitors/components/MonitorResultsList.tsx:36` and
  `src/routes/_authenticated.tsx:16` — the two consumers. Both pass only
  `className` (`h-full` / `h-svh`). Neither used Radix's `type` or
  `scrollHideDelay`, so no call site changed.

## Behavior changes

- Scrollbar visibility model changed. Radix's `type` (default `"hover"`) and
  `scrollHideDelay` (default 600ms) are gone; Base UI mounts a scrollbar only
  when the viewport is scrollable and exposes `data-hovering`, `data-scrolling`
  and `data-has-overflow-x/y` for opacity styling instead. This wrapper sets no
  opacity classes, so scrollbars are now visible whenever the content overflows
  rather than fading in on hover and out after a delay. Flagged, not patched —
  restoring the fade means adding an opacity transition keyed on
  `data-hovering`/`data-scrolling`, which is a styling decision.
- `data-state="visible" | "hidden"` is no longer emitted on the scrollbar.
  Nothing in the project selected on it.
- The horizontal scrollbar now only mounts when there is horizontal overflow
  (previously `type="hover"` governed it). With `h-svh`/`h-full` containers and
  no horizontal overflow, expect it to be absent from the DOM rather than
  present-but-hidden.

## Verify by hand

1. `_authenticated` shell: navigate to a long page (Monitors list). Vertical
   scrolling must work and the 2.5-wide thumb appear over the content, with the
   `h-svh` viewport not double-scrolling against the window.
2. Monitor detail → results list: scroll the panel, confirm the thumb tracks
   position and dragging the thumb scrolls the content.
3. Shrink the window until something overflows horizontally and confirm the
   horizontal scrollbar appears and works — this is the preserved customization.
4. Note whether always-visible scrollbars over the content look acceptable; if
   not, that is the flagged fade-on-hover change, not a regression.
