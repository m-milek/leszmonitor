# switch

2026-09-20, transformation engine (legacy `new-york` style: classified against
the stock `new-york-v4` golden, then transformed in place). Migrated 1:1 to
`@base-ui/react/switch`.

## Changed

- `src/components/ui/switch.tsx` — `import { Switch as SwitchPrimitive } from "radix-ui"`
  → `from "@base-ui/react/switch"`. Parts are unchanged (`Root`, `Thumb`).
  Props type `React.ComponentProps<typeof SwitchPrimitive.Root>` →
  `SwitchPrimitive.Root.Props` (the `size` intersection is kept); the now-unused
  `import * as React` was removed. Class rewrites per `class-mapping.md`:
  - Root `disabled:cursor-not-allowed disabled:opacity-50` →
    `data-disabled:cursor-not-allowed data-disabled:opacity-50`. Required, not
    cosmetic: Base UI's Switch Root renders a `<span>` plus a hidden `<input>`,
    so the `disabled:` pseudo-class variants were dead code.
  - Root `data-[state=checked]:bg-primary` → `data-checked:bg-primary`,
    `data-[state=unchecked]:bg-input` → `data-unchecked:bg-input`,
    `dark:data-[state=unchecked]:bg-input/80` → `dark:data-unchecked:bg-input/80`.
  - Thumb `data-[state=checked]:translate-x-[calc(100%-2px)]` →
    `data-checked:…`, `data-[state=unchecked]:translate-x-0` → `data-unchecked:…`,
    `dark:data-[state=checked]:bg-primary-foreground` → `dark:data-checked:…`,
    `dark:data-[state=unchecked]:bg-foreground` → `dark:data-unchecked:…`.
  - The project's own `data-[size=…]` and `group-data-[size=…]/switch` hooks are
    untouched (custom attribute, not a Radix state attribute), as are
    `data-slot`, `peer`, `group/switch` and every focus-visible class.
    Bare `data-checked:` / `data-disabled:` variants are valid here — Tailwind v4,
    and `src/components/ui/combobox.tsx:151` already uses `data-highlighted:`.
    Leftover scan: `grep -n "radix-ui\|@radix-ui\|data-\[state=" src/components/ui/switch.tsx` → clean.

Classification: PRISTINE against `new-york-v4`.

## Left alone

- `src/components/form/LMSwitch.tsx:15` and
  `src/features/monitors/forms/MonitorForm.tsx:139` — the two consumers. Both
  pass only `id` / `name` / `checked` / `onCheckedChange`, all unchanged names.
  `onCheckedChange` gained a second `eventDetails` argument in Base UI; both
  handlers are declared `(checked: boolean) => void`, which stays type-safe.
- `src/components/form/ErrorTooltip.tsx` — wraps the switch in a still-Radix
  `TooltipTrigger asChild`. Radix's Slot clones the child and forwards a ref;
  Base UI's Switch Root forwards refs, so the composition holds. It will be
  revisited when tooltip is migrated, not now.

## Behavior changes

- Root element changed from `<button role="switch">` to `<span>` + a hidden
  `<input type="checkbox">`. Consequence for `LMSwitch`, which passes
  `id={props.name}`: with Base UI's default `nativeButton={false}` that `id`
  now lands on the hidden input rather than the visible control. No
  `<label htmlFor>` in this codebase targets a switch id, so nothing is broken
  today — but a future label must point at the switch's `id` knowing it resolves
  to the input. `nativeButton` + `render` would restore a real `<button>`;
  flagged, not applied (the shadcn base registry keeps the span).
- The switch now submits `"on"` inside a form by default (native checkbox
  semantics) and submits nothing when off. Radix's hidden input only appeared
  inside a form. Neither consumer relies on native form submission —
  `MonitorForm` is TanStack Form and reads `checked` from state.

## Verify by hand

1. Monitors → New Monitor: toggle "Use Custom Slug". The thumb must slide the
   full width (`data-checked:translate-x-[calc(100%-2px)]`) and the track turn
   primary — this is the main check that the `data-checked:` rewrite took.
2. Tab to the switch: the 3px focus ring must appear; Space/Enter toggles it.
3. Toggle it and confirm the dependent "Slug" field enables/disables, i.e.
   `onCheckedChange` still fires with the boolean first.
4. Switch to dark mode and toggle again — thumb colour must flip
   (`dark:data-checked:bg-primary-foreground` / `dark:data-unchecked:bg-foreground`).
5. If you have a disabled switch anywhere, confirm it dims at 50% opacity
   (`data-disabled:`) — this path changed attribute.
