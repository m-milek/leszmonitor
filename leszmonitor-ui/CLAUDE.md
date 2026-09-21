# leszmonitor-ui

React 19 + Vite, TanStack Router (file-based) / Query / Form, zustand, shadcn/ui, Tailwind 4, zod.

## Commands

```bash
npx tsc --noEmit        # typecheck (`npm run build` runs tsc AFTER vite, so don't rely on it)
npm run test            # vitest
npx eslint src          # not `eslint .` — that also lints dist/
npm run build           # WARNING: wipes and overwrites ../leszmonitor-server/src/static
cd e2e && npx playwright test   # needs the server already running on :7001
```

## Structure

```
src/
  app/          main.tsx, store.ts, providers/
  routes/       thin: createFileRoute + head + loader + <XPage/> only
  features/<name>/
    <name>-api.ts   one exported namespace object, e.g. TagsApi
    types.ts        domain types
    components/ forms/ pages/ lib/ hooks/   only when there are several files
  components/
    ui/         shadcn (base-nova) — generated, don't hand-edit; style via tokens + wrappers
    common/     app-agnostic primitives (Flex, Typography, DataTable, ...)
    form/       LM* form inputs + field-state.ts
    layout/     sidebar
  lib/          api-client (authFetch), consts, cookies, jwt, shared types
```

## Conventions

- Imports use the `@/` alias and **no file extensions**.
- API: add methods to the feature's `XxxApi` object with short names (`TagsApi.getAll`, not `getAllTags`). Don't export
  loose fetch functions. Monitor results/stats live under `MonitorsApi.results` / `MonitorsApi.stats`.
- Page components receive route params as props; only the route file calls `Route.useParams()`.
- Named exports only; props interface is `{ComponentName}Props`; use `cn()` from the `cn` package for classes
  (not `@/lib/utils` — that only has `formatDate`/`formatDuration`).
- Component files `PascalCase.tsx`, everything else `kebab-case.ts`. Tests sit next to their subject.
- Auth token: go through `features/auth/lib/token.ts`. `readTokenSync` (router guard) and `readToken` (async) are
  intentionally separate.
- Never edit `src/routeTree.gen.ts` — it's generated from `src/routes/`.
- Mutations default to `throwOnError: true` (`app/main.tsx`), so a failed mutation surfaces
  instead of being swallowed. To handle one in place, opt out per call:
  `useMutation({ throwOnError: false, onError })`. There is no global error UI.

## Styling

- `src/styles.css`: the `:root` / `.dark` / first `@theme inline` blocks come from the theme
  generator — don't hand-edit them. Project-specific tokens (`--lm-status-*`, `--font-heading`)
  and `@layer base` live in the clearly marked block at the bottom of the file.
- Status colours go through `--lm-status-{up,down,pending,paused,unknown}`; never raw Tailwind
  palette classes (`bg-green-500`, `text-emerald-500`, …). `StatusDot` exports `STATUS_BG_CLASS`.
- Don't edit `src/components/ui/*`. Change appearance via variants, semantic tokens, CSS
  variables, or a wrapper component. The remaining intentional forks there are: `field.tsx`
  (`FieldContext` so `FieldLabel` derives `htmlFor`), `dropdown-menu.tsx`
  (`DropdownMenuItemIcon`), `tooltip.tsx` (`arrowClassName`, `role="tooltip"`, `z-[9999]`,
  self-wrapping `TooltipProvider`), `scroll-area.tsx` (both scrollbars), `card.tsx`
  (`ring-foreground/20`), `toast.tsx` (icon colours), `button.tsx` (`cursor-pointer`),
  `sidebar.tsx` (`SIDEBAR_WIDTH = 18rem`, overridden per-app via the `--sidebar-width` style
  prop in `routes/_authenticated.tsx`). Re-apply them after `npx shadcn add <component>`.
