# Copilot instructions

## Project summary

- Leszmonitor: lightweight homelab monitoring.
- Backend: Go service in `backend/src`.
- Frontend: web app in `frontend`.

## Workflow (Taskfile)

- Use the Taskfile as the primary entry point for developer tasks.
- Discover available tasks with:

```bash
task --list
```

- Prefer `task <name>` over invoking tools directly unless a task does not exist.

- Avoid using terminal commands - use your tools. Avoid gluing together some complicated commands - use tools.

## Coding guidelines

- Keep changes scoped and consistent with existing patterns in each area (Go backend, frontend).
- Avoid introducing new dependencies unless necessary; if you add one, update the relevant manifest.
- Add tests only when the feature is trivial to test; if not, ask before proceeding.
- Pay attention to backend security practices, especially around authentication and data handling.

## Backend guidelines

- SQLite database
- JWT auth, user data in context.Context

## Frontend guidelines

- Prefer named exports over default exports
- All components use props interface named `{ComponentName}Props`
- Use `cn()` utility from `@/lib/utils` for conditional classes

## When unsure

- Check `README.md` and `Taskfile.yml` first.
- If a task is unclear, ask before guessing.
