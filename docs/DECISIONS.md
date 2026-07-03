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

## D-009: Frontend API request contract is pinned by tests

**Date:** 2026-07-02 (issue Pollyglot#12)

**Context:** The register form 404'd in manual testing because the env base
URL and the service paths both carried `/v1` (Pollyglot#19). Nothing caught
it: component tests mock the service layer, and the backend tests don't know
what URLs the frontend uses.

**Decision:** `decks.service.test.ts` mocks the axios client and asserts the
exact path, verb, and payload of every service call. New services get the
same contract suite.

**Why:** It is the only cheap place to catch path drift (double prefixes,
renamed routes) before a human hits it in the browser.

## D-010: Study sessions run on a client-side queue snapshot

**Date:** 2026-07-02 (issue Pollyglot#13)

**Context:** The study queue endpoint returns due cards. TanStack Query
refetches on window focus, and each review changes what is due — a live
query would reorder or shrink the deck mid-session.

**Decision:** The session copies the queue into local state the first time
it loads and iterates over that snapshot; reviews post per card, and a
failed review keeps the learner on the same card with an inline error
instead of advancing.

**Why:** A learner mid-session must never see cards shuffle under them, and
a lost network blip must not silently drop a review (the rating the user
chose is the product's core data).

## D-011: Repo-wide `make check` mirrors CI exactly

**Date:** 2026-07-02

**Context:** The gates live in two toolchains (Go and bun) and CI runs both;
running them by hand meant remembering four commands in two directories.

**Decision:** A root Makefile provides `make check` (everything CI runs),
`make test`, `make dev`, and `make dev-down`. Shared frontend test helpers
(`src/lib/test-utils.tsx`: `renderWithQuery`, `mockDeck`, `mockCard`) remove
the per-file QueryClient/factory boilerplate.

**Why:** One command that equals CI keeps "all tests pass before merging"
honest locally, and shared factories keep exhaustive test suites cheap to
write.

## D-012: Startup seeder with a dev-only tier

**Date:** 2026-07-02 (issue Pollyglot#20)

**Context:** Marc wanted a seeded account to sign in without registering.
The starter shipped `migrations/seeders/rbac_seeder.sql` but nothing ever
executed it, so RBAC baseline data silently never landed.

**Decision:** `pkg/seed` runs after migrations on every boot:
`migrations/seeders/*.sql` in every environment, `migrations/seeders/dev/`
only when `-env development`. Files run in lexicographic order and must be
idempotent (`ON CONFLICT DO NOTHING`). Dev tier seeds
`demo@pollyglot.dev` / `Password123!` (fixed UUIDs, bcrypt-hashed) plus a
six-card starter deck.

**Why:** Seeds-as-SQL keeps them reviewable and idempotent-by-construction;
running on boot (not as migrations) keeps reference/dev data out of the
schema history; the dev/ split makes it impossible to ship the demo account
to production. Fixed UUIDs make the seeded data addressable from tests.

## D-013: Factory pattern for all test data

**Date:** 2026-07-02 (issue Pollyglot#21)

**Context:** Marc asked for factory-pattern test data (`UserFactory`,
`DeckFactory`, `CardFactory`) with seeded users/values. Tests were building
models with inline literals duplicated per file.

**Decision:** Go: `internal/shared/factory` with chainable builders
(`factory.Card().WithDeckID(id).WithSRS(1.9, 12, 4).Build()`) and a
`factory.Seeded` struct pinning the dev-seed UUIDs/credentials (a test
asserts they stay in sync with the SQL). TS: `UserFactory` / `DeckFactory` /
`CardFactory` in `src/lib/test-utils.tsx` with `build`/`buildList` and
sequence-numbered defaults, plus `SeededUser`. All suites refactored; new
tests must use factories.

**Why:** One place defines what a "normal" user/deck/card looks like, so
schema changes touch one file instead of every test; sequence-numbered
defaults prevent accidental cross-test identity collisions; pinning the
seeded fixtures in code keeps manual-testing credentials and automated tests
from drifting apart. Conventions recorded in CLAUDE.md so future sessions
follow them.

## D-014: Translation ships with a built-in dictionary provider

**Date:** 2026-07-02 (issue Pollyglot#14)

**Context:** D-007 mandates deterministic defaults behind provider
interfaces. The translate feature needs to do *something* real without ML
infra or API keys.

**Decision:** `Translator` interface selected by `TRANSLATOR_PROVIDER`
(default `dictionary`). The dictionary provider embeds a small bidirectional
word list (Japanese/Spanish/French/German ↔ English) that includes every
word in the dev-seeded starter deck; unknown input returns 422 "no
translation available" (provider outages would be 502). The UI treats 422 as
a friendly message, not an error state.

**Why:** A lookup that really translates the demo content makes the feature
honest end-to-end (demo account → starter deck → translate → save back as a
card) while keeping tests hermetic. Distinguishing 422 (no data) from 502
(provider broken) means the future ML/LLM provider drops in without any
handler or UI changes.

## D-015: Conversation tutor is the scripted Socratic persona; exchanges are atomic

**Date:** 2026-07-02 (issue Pollyglot#15)

**Context:** The original repo shipped a Socratic-tutor prompt ("never give
answers, always end with a question"). No LLM is available in this
environment, and a chat that half-persists on failure corrupts history.

**Decision:** `TutorProvider` interface with a deterministic
`SocraticTutor`: quotes the learner's words, cycles five distinct probes by
turn count, always ends with a question (tests enforce all three
properties). `POST /conversations/{id}/messages` asks the provider *before*
persisting anything, then stores the user and tutor messages together and
returns the full exchange; the UI appends the exchange into the query cache
instead of refetching.

**Why:** The scripted persona preserves the product's pedagogy honestly
without faking an AI, and gives an LLM provider a behavioral contract to
meet (the tests document it). Provider-before-persist means a provider
outage is a clean 502 with no orphaned user message to dedupe later;
cache-append keeps the chat from flickering on every turn.

## D-016: Stats aggregate in SQL by day; streak logic is a pure function

**Date:** 2026-07-02 (issue Pollyglot#16)

**Context:** Progress needs reviews-per-day (chart), counts, and a streak.
Streak semantics have edge cases (is a streak "broken" before you study
today?), and naive row-fetching would pull every review to Go.

**Decision:** The repository does one `GROUP BY` day query (a year back, so
long streaks survive) plus two counts; `Streak()` and `FillDays()` are pure
functions — streak counts back from today, or from yesterday when today is
still unreviewed, so a live streak never displays as broken; future-dated
noise is ignored. The chart series is always exactly 30 zero-filled days.
Chart colors were validated with the dataviz six-checks script per surface
(#059669 light / #0ea371 dark), single hue for a single series, with an
sr-only table of the same data.

**Why:** SQL aggregation keeps the payload tiny; pure functions made the
ten streak edge cases table-testable; the "yesterday keeps it alive" rule
matches learner expectations (Duolingo-style) instead of punishing morning
visits; validated color + a data table keep the chart accessible in both
themes.

## D-017: i18n covers navigation now; page copy stays English pending a dedicated sweep

**Date:** 2026-07-02 (issue Pollyglot#17)

**Context:** The polish issue asked for next-intl keys on all new pages.
Surveying the codebase: the starter itself only uses `useTranslations` in
two components — every other page hardcodes English (and the logout dialog
hardcoded Indonesian). A full retrofit of nine product pages would churn
every component test late in the build-out while exceeding the starter's
own i18n adoption.

**Decision:** Internationalize the shared chrome now — sidebar labels are
locale keys (`nav.*`) resolved at render, the logout dialog uses `auth.*`
keys, EN/ID catalogs extended — and enforce catalog health with a
locale-parity test (both files must define exactly the same key paths; it
already caught one pre-existing drift). Page-level copy remains English;
a dedicated i18n sweep is deferred to the next backlog-planning session
with Marc.

**Why:** Navigation is the highest-visibility string surface and now
switches languages correctly; the parity test turns missing translations
into CI failures instead of raw keys in production; and deferring the
page sweep honestly (documented here) beats a rushed half-translation
that would still fail the "all pages" bar.

## D-018: Neumorphic surface system with per-theme shadow pairs

**Date:** 2026-07-02 (issue Pollyglot#27, requested by Marc)

**Context:** Marc asked for a neumorphic style. The original Pollyglot repo
already used soft UI (`#e0e5ec` surfaces, dual-shadow `.neumorphic-shadow`
classes), so this restores the product's original design language.
Neumorphism's known risk is contrast collapse (low-contrast surfaces,
borderless controls).

**Decision:** Four utility classes in globals.css — `neu-card`,
`neu-card-sm`, `neu-inset`, `neu-btn` (+ `neu-interactive` hover/press) —
driven by two per-theme shadow variables. Light theme uses the classic
`#e0e5ec` family with white/`#b8bec7` shadows (matching the original repo);
dark theme lifts the surface to a midtone charcoal with its own tuned pair
(derived from the original's oklch values), never an automatic flip.
Guardrails: text keeps the standard foreground tokens (AA on both
surfaces), emerald stays the action accent, focus rings untouched, raised
cards press inset on activation for tactile feedback. Applied to the
Pollyglot product surface only; starter admin/auth pages unchanged.

**Why:** Per-theme tuning is what makes dark neumorphism legible (a
flipped light shadow pair reads as dirt on dark surfaces); scoping to the
product surface keeps the diff reviewable; and empty states as recessed
wells vs. content as raised cards gives the soft-UI hierarchy an actual
meaning instead of decoration.

## D-019: Runtime upgrades are staged, tagged, and data-backed

**Date:** 2026-07-02 (issue Pollyglot#30, requested by Marc)

**Context:** Go 1.24 passed EOL in February 2026; Postgres 16, Redis 7,
and Next 16.1 all had newer stables. Marc required a backup at every stage
so any single upgrade can be reverted.

**Decision:** One branch, one atomic commit per component, each verified
live before the next, with three layers of rollback: (1) a git tag before
every stage (`pre-upgrade-{baseline,go,redis,postgres,nextjs}`); (2)
`pg_dump` SQL + Redis RDB snapshots in gitignored `backups/`; (3) Postgres
migrated by dump→restore into a **new** volume (`postgres-data-18`,
mounted at `/var/lib/postgresql` per the PG18 image) while the PG16
volume stays untouched — reverting the compose file boots the old data
unchanged. Restore was verified by per-table row-count comparison
(0 errors, identical counts). Landed: Go 1.26.4, Redis 8.8.0,
Postgres 18.4, Next.js 16.2.10.

**Why:** Tags make code rollback one command; the parallel Postgres
volume makes data rollback instant rather than a restore drill; per-stage
verification means a failure pins to one component instead of one big-bang
upgrade to bisect.

## D-020: Due counts ride the deck list; the daily goal lives on the user

**Date:** 2026-07-02 (issue Pollyglot#22)

**Context:** Learners need to see where work is waiting (due cards) and
have a target (daily goal). Options were a separate due-counts endpoint vs
enriching the existing deck payload, and a settings module vs a column.

**Decision:** `DeckResponse` gains `due_count` computed alongside
`card_count` in the deck list/detail queries — no new endpoint, badges
render wherever decks render. The goal is a `users.daily_goal` column
(default 20, DB CHECK > 0) read/written through the stats module
(`GET /v1/stats` includes it; `PUT /v1/stats/goal`, validated 1–500 with a
pointer DTO so 0 is rejected not ignored). The stats page gets a
progressbar (proper ARIA value semantics) with a met-state celebration;
"met" is `reviews_today >= goal`, so raising the goal mid-day can honestly
un-meet it.

**Why:** The deck list is already the surface where "what should I study"
gets decided — a second request per deck would just be latency; a full
settings module for one integer is ceremony, and the stats module already
owns the review-count context that gives the goal meaning.

## D-021: Reverse is card duplication; cloze is a card type; pronunciation is browser TTS

**Date:** 2026-07-02 (issue Pollyglot#23)

**Context:** "Richer cards" could mean per-direction SRS state on one card
(Anki's note/card split), a template system, or server-side audio.

**Decision:** Three deliberately-small mechanisms. **Reverse** is a
create-time option that persists a second, independent basic card with
mirrored front/back and fresh SRS state — no schema change, each direction
schedules on its own; rejected for cloze. **Cloze** is a `card_type`
('basic'|'cloze', DB CHECK) with Anki-syntax `{{c1::…}}` markers parsed by
twin pure-function libraries (`backend/pkg/cloze` ↔ `frontend/src/lib/
cloze.ts`, same table tests both sides); creation requires ≥1 well-formed
deletion. **Pronunciation** is the browser SpeechSynthesis API
(feature-detected, language-name → BCP-47 mapping) — zero backend, and the
upcoming ElevenLabs provider (Pollyglot#28) slots in as the premium path
with this as the fallback.

**Why:** Duplication gives reverse cards correct independent scheduling
for free, exactly what a per-direction state machine would have bought at
10× the complexity; a type column keeps cloze additive (existing cards
untouched); and mirrored parser implementations with identical test tables
keep the one format that crosses the API boundary honest on both sides.

## D-022: Server speech is optional; the client always has a voice

**Date:** 2026-07-02 (issue Pollyglot#28, requested by Marc)

**Context:** Marc wants ElevenLabs voices for conversation practice. API
keys may be absent (dev machines, CI), and audio must never be the reason
a page breaks.

**Decision:** `internal/speech` follows D-007: a `Provider` interface with
an ElevenLabs implementation (multilingual model, key in the `xi-api-key`
header server-side only, stub-server tests — never the network).
`NewProvider` returns nil when unconfigured and `POST /v1/speech` answers
503; the client's `speakWithFallback` tries server audio first and falls
back to browser SpeechSynthesis on *any* failure, so the play button on
tutor bubbles always works. Study-card pronunciation stays on browser TTS
(single words don't justify API spend).

**Why:** The 503-plus-fallback contract makes the premium voice a pure
enhancement — zero configuration works everywhere, one env var upgrades
it; and keeping the key behind our API means it never ships to browsers.

## D-023: Google Translate slots behind the existing Translator interface

**Date:** 2026-07-02 (issue Pollyglot#29, requested by Marc)

**Context:** Marc wants real translations. D-014 already split 422 (no
translation) from 502 (provider broken) precisely so a real provider could
drop in without touching handlers or UI.

**Decision:** `GoogleTranslator` (Cloud Translation v2, JSON body,
`format:text` to avoid HTML entities, key server-side) selected by
`TRANSLATOR_PROVIDER=google` + `GOOGLE_TRANSLATE_API_KEY`; keyless or
unknown values fall back to the dictionary with a startup warning. Human
language names map to ISO-639-1 through a pure `LanguageCode` function;
unmappable languages return `ErrNoTranslation` *before any request leaves
the process*. Tests run against an httptest stub pinning the request
shape and error mapping.

**Why:** The name→code mapping is the one place user vocabulary meets the
Google API contract, so it's pure and table-tested; refusing unmappable
languages locally keeps quota for real requests and keeps the 422/502
split honest; and zero changes were needed in the service, handler, or UI
— which is exactly what D-014 promised.

## D-024: Import/export is CSV/TSV with per-row error reporting

**Date:** 2026-07-02 (issue Pollyglot#24)

**Context:** Bulk vocabulary needs a way in and out. Anki compatibility
matters (TSV), and imports of hand-edited files will contain bad rows.

**Decision:** Pure serializer/parser functions in the decks module
(stdlib encoding/csv, comma or tab): header row `front,back,card_type`
with the type column optional on import; malformed rows are skipped and
reported with 1-indexed line numbers (header counted, so numbers match
the user's editor) while good rows import; whole-file failures only for
unknown formats and the 1000-row cap. Imported cards always start with
fresh SRS state. A round-trip property test pins export→import fidelity.
Live verification caught that the router's global
`AllowContentType("application/json")` middleware 415'd multipart
uploads — now allows `multipart/form-data`.

**Why:** Skip-and-report beats all-or-nothing for hand-edited files (one
typo shouldn't block 300 words) while the row cap bounds abuse; fresh SRS
state on import is honest (we don't know the source app's scheduler); and
the 415 catch is the standing argument for live-verifying every feature
against the real stack.

## D-025: Sharing is a capability code; cloning is a fresh copy

**Date:** 2026-07-02 (issue Pollyglot#25)

**Context:** Deck sharing could mean live collaboration, follower decks,
or copies. Access control options ranged from public flags to ACLs.

**Decision:** A nullable unique `share_code` on the deck (10 chars from
an unambiguous alphabet, crypto/rand, collision-retried against the
unique index) acts as a capability: any *authenticated* user with the
link can preview (name, languages, count, five sample cards) and clone.
Sharing is idempotent, revocable (code cleared), and codes are not
enumerable. Clones are independent copies — cloner's ownership, fresh
SRS state, unshared — with no backlink to the source.

**Why:** Capability links are the simplest sharing model that matches
the product's social reality (send a friend your deck); requiring auth
keeps content behind accounts without an ACL system; fresh-copy
semantics avoid the consistency swamp of live-shared decks (what happens
to reviews when the owner edits?) — a clone is yours, full stop.
