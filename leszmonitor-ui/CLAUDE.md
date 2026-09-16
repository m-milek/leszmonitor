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
    ui/         shadcn — generated, don't hand-edit
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
- Named exports only; props interface is `{ComponentName}Props`; use `cn()` from `@/lib/utils` for classes.
- Component files `PascalCase.tsx`, everything else `kebab-case.ts`. Tests sit next to their subject.
- Auth token: go through `features/auth/lib/token.ts`. `readTokenSync` (router guard) and `readToken` (async) are
  intentionally separate.
- Never edit `src/routeTree.gen.ts` — it's generated from `src/routes/`.
