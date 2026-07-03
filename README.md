<p align="center">
  <img src="frontend/public/pollyglot.svg" alt="Pollyglot" width="120" />
</p>

<h1 align="center">Pollyglot 🦜</h1>

<p align="center">
  Learn languages the way memory works — spaced-repetition flashcards,
  translation you can keep, and a tutor that asks before it answers.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Next.js-16-000000?style=flat&logo=nextdotjs&logoColor=white" />
  <img src="https://img.shields.io/badge/Tailwind_CSS-v4-38BDF8?style=flat&logo=tailwindcss&logoColor=white" />
  <img src="https://img.shields.io/badge/PostgreSQL-18-4169E1?style=flat&logo=postgresql&logoColor=white" />
  <img src="https://img.shields.io/badge/Redis-8-DC382D?style=flat&logo=redis&logoColor=white" />
</p>

Built by [Cometra](https://github.com/CometraLLC) on the GNS starter kit
(Go + Next.js + shadcn/ui: JWT auth, RBAC, i18n, clean architecture).

## ✨ Features

- **Decks & flashcards** — per-user decks for each language you learn
- **Spaced repetition** — SM-2 scheduling; rate a card Forgot → Got it! and it comes back just before you'd lose it
- **Translate & keep it** — dictionary-backed (pluggable) translation with one-click save-to-deck
- **Conversation practice** — a Socratic tutor that quotes your words and always ends with a question (pluggable provider, LLM-ready)
- **Progress** — day streak, reviews per day, unique words, 30-day chart
- **Platform** (from GNS) — JWT auth + Google OAuth, RBAC, i18n (EN/ID), dark mode, migrations, rate limiting, security headers

## 🏗️ Tech Stack

| Layer        | Technology                                       |
| ------------ | ------------------------------------------------ |
| **Backend**  | Go, Chi Router, uber/dig, PostgreSQL, Redis      |
| **Frontend** | Next.js 16, React 19, Tailwind CSS v4, shadcn/ui |
| **Auth**     | JWT (access + refresh tokens), Google OAuth      |
| **State**    | TanStack Query, Zustand                          |
| **i18n**     | next-intl (EN, ID)                               |
| **DevOps**   | Docker, Docker Compose, Makefile                 |

## 📊 Feature status

All product pages live under `/pollyglot` after signing in:

| Feature | Status | Where |
| ------- | ------ | ----- |
| Decks & cards CRUD | ✅ | `/pollyglot/decks`, API `/v1/decks` |
| Spaced-repetition study (SM-2, cloze + reverse cards, TTS) | ✅ | `/pollyglot/study`, API `/v1/cards/{id}/review` + `/v1/decks/{id}/queue` |
| Translate (dictionary provider, save-to-deck) | ✅ | `/pollyglot/translate`, API `/v1/translate` |
| Conversation practice (Socratic tutor) | ✅ | `/pollyglot/conversation`, API `/v1/conversations` |
| Progress & stats (streak, 30-day chart, daily goal) | ✅ | `/pollyglot/stats`, API `/v1/stats` |

Every engineering decision is logged with rationale in
[`docs/DECISIONS.md`](docs/DECISIONS.md). Work is issue-driven on
[CometraLLC/Pollyglot](https://github.com/CometraLLC/Pollyglot/issues); each
issue lands via a test-first PR that must be green before merge.

### Development commands (repo root)

```bash
make check   # everything CI runs: backend lint+tests, frontend lint+tests+build
make test    # all tests only
make dev     # docker stack (API :8080) + Next dev server (:3000)
make dev-down
```

### Optional provider keys (backend/.env)

| Variable | Effect |
| -------- | ------ |
| `SPEECH_PROVIDER=elevenlabs` + `ELEVENLABS_API_KEY` | Tutor messages play with a natural ElevenLabs voice; without it the UI falls back to the browser's built-in speech. |
| `TRANSLATOR_PROVIDER=google` + `GOOGLE_TRANSLATE_API_KEY` | Translate uses the Google Cloud Translation API; without it the built-in dictionary handles the seeded vocabulary. |

### Seeded dev account

In development the API seeds an account on startup so you can sign in
without registering:

| Email | Password |
| ----- | -------- |
| `demo@pollyglot.dev` | `Password123!` |

It comes with a "Japanese Basics" starter deck (6 cards, all due) so decks
and study have content immediately. Seed files live in
`backend/migrations/seeders/` (every environment) and
`backend/migrations/seeders/dev/` (development only); all seeds are
idempotent.

## 🚀 Quick Start

### Prerequisites

- Go 1.26+
- Node.js 18+ or Bun
- Docker & Docker Compose

### 1. Clone

```bash
git clone https://github.com/yogameleniawan/gns.git
cd gns
```

### 2. Backend

```bash
cd backend
cp .env.example .env
cp config/config.template.yaml config/config.development.yaml

# Option A: Docker (recommended)
make dev

# Option B: Local (requires PostgreSQL & Redis)
go mod download
make init
make run
```

### 3. Frontend

```bash
cd frontend
bun install    # or: npm install
bun dev        # or: npm run dev
```

### 4. Open

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080

### Default Admin Account

| Field    | Value           |
| -------- | --------------- |
| Email    | `admin@gns.com` |
| Password | `admin123`     |

## 📁 Project Structure

```
gns/
├── backend/
│   ├── cmd/api/          # Entry point (main, config, migration, server)
│   ├── config/           # YAML config per environment
│   ├── container/        # Dependency injection (uber/dig)
│   ├── internal/         # Business modules
│   │   ├── auth/         # Authentication (dto, repo, service, handler)
│   │   └── rbac/         # Roles & permissions
│   ├── migrations/       # SQL migrations & seeders
│   └── pkg/              # Shared packages (middleware, router, utils, etc.)
│
├── frontend/
│   ├── app/              # Next.js App Router pages
│   ├── src/
│   │   ├── domain/       # Types & interfaces
│   │   ├── application/  # Hooks & state management
│   │   ├── infrastructure/ # API clients & stores
│   │   └── presentation/ # Components & pages
│   └── locales/          # i18n translations
```

## 🔧 Backend Commands

```bash
make dev                    # Start dev environment (Docker)
make dev-down               # Stop dev containers
make run                    # Run locally
make build                  # Build binary
make create-migration name=X  # Create migration
make migrate-up             # Run migrations
make migrate-down           # Rollback migrations
make db-shell               # PostgreSQL shell
```

## 📖 Adding a New Module

Each backend module follows a 4-file pattern inside `internal/`:

```
internal/your_module/
├── dto.go          # Request/response structs
├── repository.go   # Database queries
├── service.go      # Business logic
└── handler.go      # HTTP handlers
```

Then wire it in `container/container.go` and add routes in `pkg/router/router.go`.

See the [Documentation page](/docs) for detailed guides.

## 🌐 API Routes

All routes are prefixed with `/v1`.

| Method | Path                | Description                 |
| ------ | ------------------- | --------------------------- |
| POST   | `/auth/register`    | Register                    |
| POST   | `/auth/login`       | Login                       |
| POST   | `/auth/refresh`     | Refresh token               |
| GET    | `/auth/profile`     | Get profile (🔒)            |
| GET    | `/users`            | List users (🔒 Admin)       |
| POST   | `/users`            | Create user (🔒 Admin)      |
| GET    | `/rbac/roles`       | List roles (🔒 Admin)       |
| GET    | `/rbac/permissions` | List permissions (🔒 Admin) |

🔒 = Requires authentication

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
