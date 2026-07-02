import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render } from '@testing-library/react'
import type { ReactElement } from 'react'
import type { Card, Deck } from '@/src/domain/services/decks.service'

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

/** Deck factory with sensible defaults; override what the test cares about. */
export function mockDeck(overrides: Partial<Deck> = {}): Deck {
	return {
		id: 'deck-1',
		name: 'Japanese Basics',
		source_lang: 'Japanese',
		target_lang: 'English',
		card_count: 3,
		created_at: '2026-07-01T00:00:00Z',
		updated_at: '2026-07-01T00:00:00Z',
		...overrides,
	}
}

/** Card factory with sensible defaults; override what the test cares about. */
export function mockCard(overrides: Partial<Card> = {}): Card {
	return {
		id: 'card-1',
		deck_id: 'deck-1',
		front: 'こんにちは',
		back: 'hello',
		ease_factor: 2.5,
		interval_days: 0,
		repetitions: 0,
		due_at: '2026-07-01T00:00:00Z',
		created_at: '2026-07-01T00:00:00Z',
		updated_at: '2026-07-01T00:00:00Z',
		...overrides,
	}
}
