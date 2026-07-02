# Pollyglot — Decision Log

A running log of engineering decisions made while building Pollyglot on the GNS
starter kit. Each entry records the context, the decision, and why. Newest
entries are appended to the bottom. Follows on from the ADRs in the original
`CometraLLC/Pollyglot` repo (`docs/adr/0001-Platform.md`, `0002-Languages.md`).

---

## D-001: Backlog lives as GitHub issues on `CometraLLC/Pollyglot`

**Date:** 2026-07-02

**Context:** The build-out is issue-driven ("get all the issues from GitHub and
iterate"), but no issues existed anywhere: `pollyglot-gns` has issues disabled
(and enabling them is a repo-settings change not authorized here), and the
sibling `CometraLLC/Pollyglot` repo had zero issues in any state. A GitHub
Project board may exist but the local `gh` token lacks the `read:project`
scope (`gh auth refresh -s read:project` would grant it).

**Decision:** Author the backlog as GitHub issues on `CometraLLC/Pollyglot`
(the product repo, which already has issues enabled) and reference them from
branches/commits in `pollyglot-gns` as `CometraLLC/Pollyglot#N`.

**Why:** Keeps an auditable, issue-driven workflow without changing repo
settings. Issues can be migrated to `pollyglot-gns` later if issues get
enabled there.

## D-002: Product scope is derived from the original Pollyglot repo

**Date:** 2026-07-02

**Context:** No PRD exists. The original `CometraLLC/Pollyglot` repo contains
the product's DNA: a study page with a flashcard component (flip + five-level
rating: Forgot / Difficult / Okay / Almost / Got it!), "Cards Flipped" and
"Unique Words" counters, Jest tests for the flashcard flow, a `ml/` PyTorch
service stub, a Socratic-tutor prompt, and ADRs. The in-progress migration in
`pollyglot-gns` adds `app/pollyglot/` with a placeholder describing Pollyglot
as "a tool that helps you translate text between languages" and an empty
`conversation/` directory.

**Decision:** Build Pollyglot as a language-learning app with these MVP
features on top of the GNS starter (Go API + Next.js + auth/RBAC):

1. Decks and flashcards (CRUD, per-user)
2. Spaced-repetition study sessions (five-level rating, due-card queue)
3. Translation tool (pluggable provider)
4. Conversation practice (chat UI, pluggable tutor provider)
5. Progress stats (reviews, unique words, streak)

**Why:** Every feature traces to an artifact the team already built or
scaffolded; nothing is invented from thin air.

## D-003: Workflow — feature branch per issue, TDD, tests gate merges to main

**Date:** 2026-07-02

**Context:** Explicit user direction: TDD with heavy testing, feature
branches, auto-push allowed, and all tests must pass before merging into main
for each issue.

**Decision:** For each issue: branch `feat/<issue>-<slug>` (or `fix/`,
`chore/`) → write failing tests first → implement → run the full test suite →
merge to `main` via PR only when green → push. PR bodies link the issue.

**Why:** Matches the requested process exactly.

**Amended 2026-07-02:** Commits are atomic — one logical change per commit —
per explicit user direction, and PRs are merged with merge commits (not
squash) so that atomic history survives on `main`. The first two PRs (#1,
#2) predate this and were squashed.

## D-004: Frontend testing stack — Vitest + Testing Library

**Date:** 2026-07-02

**Context:** The GNS frontend ships with no test runner. The original
Pollyglot app used Jest + @testing-library/react. The stack here is Next.js
16 / React 19 / Tailwind v4, with bun as the package manager.

**Decision:** Use Vitest with @testing-library/react and jsdom. Port the
original flashcard test intent to Vitest as the component is rebuilt.

**Why:** Vitest is ESM-native and needs far less configuration than Jest with
React 19/Next 16, runs fast under bun, and Testing Library assertions carry
over almost verbatim from the original Jest suite.

## D-005: Backend testing stack — stdlib `testing` + testify + httptest; fakes over DB mocks

**Date:** 2026-07-02

**Context:** The Go backend (chi + GORM + dig) has zero tests and no Makefile
test target. Modules follow handler/service/repository/dto.

**Decision:** Test with the stdlib `testing` package plus testify assertions.
Services are tested against hand-written fakes of repository interfaces;
handlers via `httptest`; business logic (e.g. the SRS scheduler) is kept in
pure functions and tested exhaustively table-driven. Repositories stay thin
(GORM calls only) and are not unit-tested against a mocked SQL layer.

**Why:** Fakes keep tests fast and behavior-focused; sqlmock-style tests
mostly re-assert GORM's SQL generation, which is low-value and brittle.
SQLite-in-memory was rejected because dialect differences from Postgres give
false confidence.

## D-006: Spaced repetition algorithm — SM-2

**Date:** 2026-07-02

**Context:** The original flashcard UI already has a five-level rating
(Forgot=0 … Got it!=4). We need a scheduler to decide when a card is next due.

**Decision:** Implement SM-2 (SuperMemo-2) with the 0–4 rating mapped onto
SM-2's 0–5 quality scale (rating+1 for passes; 0–2 treated as lapses), stored
per card: ease factor, interval days, repetitions, due date.

**Why:** SM-2 is simple, well-documented, battle-tested (Anki's ancestor),
implementable as a pure function — ideal for exhaustive TDD — and fits the
existing five-button UI without redesign.

## D-007: External AI/translation services sit behind interfaces with deterministic defaults

**Date:** 2026-07-02

**Context:** Translation and conversation practice ultimately want an ML/LLM
backend (the original repo stubs a PyTorch service). No API keys or ML infra
exist in this environment, and tests must be deterministic.

**Decision:** Define Go interfaces (`Translator`, `TutorProvider`) selected by
config/env. Ship deterministic default implementations (dictionary-backed
translator, scripted tutor) so the app works offline and tests never hit the
network. Real providers (ML service, LLM API) plug in later without touching
callers.

**Why:** Keeps TDD deterministic, avoids secrets in the repo, and preserves
the original architecture's intent (separate ML service) as a swappable
implementation detail.

## D-008: Restore the GNS shell at the app root; Pollyglot pages become the product surface

**Date:** 2026-07-02

**Context:** Mid-migration, the starter's entire `app/` tree (including the
required root `layout.tsx`, `providers.tsx`, `globals.css`) was moved into
`app/old/`, which breaks the build (Next.js requires a root layout) and turns
every starter page into an `/old/*` route. `app/pollyglot/` holds a
placeholder page.

**Decision:** Move the shell (layout, providers, globals, favicon) and the
functional starter pages (auth, admin, settings, docs) back to their original
`app/` locations, rebrand the shell as Pollyglot, make `/` the Pollyglot
landing page, and delete `app/old/` and the demo pages.

**Why:** The auth/admin/settings pages are wired to the Go backend and are
product-relevant (Pollyglot needs accounts and roles); regenerating them later
would be waste. The demo pages are starter-kit marketing with no product
value. Everything deleted stays recoverable in git history.
