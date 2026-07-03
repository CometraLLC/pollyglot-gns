import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render } from '@testing-library/react'
import type { ReactElement } from 'react'
import type { Card, Deck } from '@/src/domain/services/decks.service'
import type { Conversation, Message } from '@/src/domain/services/conversation.service'
import type { User } from '@/src/domain/services/auth.service'

/**
 * Render a component inside a fresh TanStack Query client with retries
 * disabled, so tests exercise error paths deterministically.
 */
export function renderWithQuery(ui: ReactElement) {
	const client = new QueryClient({
		defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
	})
	return render(<QueryClientProvider client={client}>{ui}</QueryClientProvider>)
}

/**
 * The development seed account (backend/migrations/seeders/dev) — use these
 * constants when a test or manual session needs a known real account.
 */
export const SeededUser = {
	id: 'a0000000-0000-4000-8000-000000000001',
	deckId: 'b0000000-0000-4000-8000-000000000001',
	email: 'demo@pollyglot.dev',
	password: 'Password123!',
} as const

interface Factory<T> {
	build(overrides?: Partial<T>): T
	buildList(count: number, overrides?: Partial<T>): T[]
}

function createFactory<T>(defaults: (index: number) => T): Factory<T> {
	let sequence = 0
	return {
		build: (overrides = {}) => ({ ...defaults(++sequence), ...overrides }),
		buildList: (count, overrides = {}) =>
			Array.from({ length: count }, () => ({ ...defaults(++sequence), ...overrides })),
	}
}

export const UserFactory = createFactory<User>((i) => ({
	id: `user-${i}`,
	email: `user-${i}@pollyglot.dev`,
	name: `Test User ${i}`,
	is_oauth: false,
	is_active: true,
	email_verified: true,
	created_at: '2026-07-01T00:00:00Z',
}))

export const DeckFactory = createFactory<Deck>((i) => ({
	id: `deck-${i}`,
	name: 'Japanese Basics',
	source_lang: 'Japanese',
	target_lang: 'English',
	card_count: 3,
	due_count: 0,
	created_at: '2026-07-01T00:00:00Z',
	updated_at: '2026-07-01T00:00:00Z',
}))

export const CardFactory = createFactory<Card>((i) => ({
	id: `card-${i}`,
	deck_id: 'deck-1',
	front: 'こんにちは',
	back: 'hello',
	card_type: 'basic',
	ease_factor: 2.5,
	interval_days: 0,
	repetitions: 0,
	due_at: '2026-07-01T00:00:00Z',
	created_at: '2026-07-01T00:00:00Z',
	updated_at: '2026-07-01T00:00:00Z',
}))

export const ConversationFactory = createFactory<Conversation>((i) => ({
	id: `conv-${i}`,
	title: 'Practice Japanese',
	language: 'Japanese',
	created_at: '2026-07-01T00:00:00Z',
	updated_at: '2026-07-01T00:00:00Z',
}))

export const MessageFactory = createFactory<Message>((i) => ({
	id: `msg-${i}`,
	role: 'tutor',
	content: 'What do you already know?',
	created_at: '2026-07-01T00:00:00Z',
}))
