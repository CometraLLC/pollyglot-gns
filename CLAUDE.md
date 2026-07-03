# CLAUDE.md

Pollyglot — a language-learning app (spaced-repetition flashcards, translation,
conversation practice) built on the GNS starter (Go + Next.js + shadcn/ui).
Work is issue-driven: issues live on **CometraLLC/Pollyglot** (this repo has
issues disabled); PRs merge into `main` here.

## Commands

```bash
make check           # everything CI runs: backend lint+tests, frontend lint+tests+build
make test            # all tests only
make dev / dev-down  # docker stack (API :8080, Postgres, Redis) + Next dev (:3000)
cd backend && make test-pkg pkg=./pkg/srs   # one Go package, verbose
cd frontend && bun run test:watch           # vitest watch mode
```

Dev sign-in without registering: `demo@pollyglot.dev` / `Password123!`
(seeded with a starter deck; see `backend/migrations/seeders/dev/`).

## Architecture

- **backend/** Go, chi router, GORM/Postgres, Redis, uber/dig DI.
  Modules in `internal/<name>/` follow **dto / repository / service / handler**;
  wire new modules in `container/container.go` and `pkg/router/router.go`
  (handlers expose `RegisterRoutes(r chi.Router, h Handler)`).
  Migrations in `migrations/` (golang-migrate, run at startup); idempotent SQL
  seeds in `migrations/seeders/` (all envs) and `migrations/seeders/dev/`
  (development only). Pure algorithms live in `pkg/` (e.g. `pkg/srs` = SM-2).
  Auth: JWT; handlers read the user via `middleware.GetUserContext(ctx)` —
  never raw context keys.
- **frontend/** Next.js 16 App Router, React 19, Tailwind v4, bun.
  Clean architecture under `src/`: `domain/services/*.service.ts` (axios API
  clients), `application/hooks/` (TanStack Query), `presentation/components/`
  (pages + shadcn ui). Route files in `app/` stay thin (ProtectedRoute +
  MainLayout + page component). API base URL is **versionless**; every service
  path carries `/v1` (see decision D-009).

## Workflow (user-mandated — do not drift)

1. Every work item is a **formal GitHub issue** on CometraLLC/Pollyglot with
   acceptance criteria and a test plan; author new ones as needed.
2. **TDD, exhaustively**: failing tests first; cover error/edge/loading/empty
   paths, not just happy paths.
3. Feature branch per issue (`feat/<n>-slug`, `fix/<n>-slug`); **atomic
   commits**; PRs merge with **merge commits** (never squash).
4. **CI must pass before merge** (`gh pr checks --watch`); `make check`
   locally mirrors CI exactly.
5. Verify features **live** against the Docker stack before shipping.
6. Before closing an issue: update `docs/DECISIONS.md` (D-###, with why) and
   the README; close with a comment linking the PR and the decisions.
7. When the backlog is empty, ask Marc — don't invent scope.

## Testing conventions

- **Factories, not literals** (both stacks):
  - Go: `internal/shared/factory` — `factory.User()`, `factory.Deck()`,
    `factory.Card()`, `factory.Review()` chainable builders (`WithX(...)`,
    `.Build()`); `factory.Seeded` holds the dev-seed UUIDs/credentials.
  - TS: `src/lib/test-utils.tsx` — `UserFactory` / `DeckFactory` /
    `CardFactory` with `.build(overrides)` / `.buildList(n)`, `SeededUser`,
    and `renderWithQuery(ui)` for TanStack Query tests.
- Go services are tested against **hand-written fakes** of the repository
  interface (no SQL mocks); handlers via `httptest` with a fake service and a
  stub auth middleware. Repositories stay thin GORM calls.
- Every frontend service gets a **request-contract test** (mock the axios
  client, assert exact path/verb/payload) — this is what catches URL drift.
- Component tests mock the service module (`vi.mock`) and use real
  TanStack Query via `renderWithQuery`.

## Gotchas

- `environtment` (sic) is the env global in `backend/cmd/api/main.go`.
- bcrypt cost is 12 — keep password-hash test counts low (each hash ~250ms).
- air (docker dev) rebuilds only on `.go` changes; `docker restart gns-dev`
  to re-run seeds after editing SQL.
- Locale files: `frontend/locales/{en,id}.json` (next-intl). Sidebar labels
  are `nav.*` keys resolved in `SidebarItem`; a locale-parity test
  (`src/lib/locales.test.ts`) fails if EN/ID key sets drift — add every new
  key to both files. Page copy is English for now (see D-017).
- **Styling is neumorphic** (D-018): new product UI uses `neu-card` /
  `neu-card-sm` (raised), `neu-inset` (recessed wells/empty states),
  `neu-btn` (pill buttons), `neu-interactive` (hover lift, press inset)
  from globals.css — not `border bg-card shadow-*`. Text keeps standard
  foreground tokens; emerald remains the accent.
