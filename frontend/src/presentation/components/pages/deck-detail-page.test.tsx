import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { Card, Deck } from '@/src/domain/services/decks.service'
import { DeckDetailPage } from './deck-detail-page'

vi.mock('@/src/domain/services/decks.service', () => ({
	decksService: {
		getDeck: vi.fn(),
		listCards: vi.fn(),
		createCard: vi.fn(),
		updateCard: vi.fn(),
		deleteCard: vi.fn(),
	},
}))

import { decksService } from '@/src/domain/services/decks.service'

const mocked = vi.mocked(decksService)

const deck: Deck = {
	id: 'deck-1',
	name: 'Japanese Basics',
	source_lang: 'Japanese',
	target_lang: 'English',
	card_count: 2,
	created_at: '2026-07-01T00:00:00Z',
	updated_at: '2026-07-01T00:00:00Z',
}

function card(overrides: Partial<Card> = {}): Card {
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

function renderPage() {
	const client = new QueryClient({
		defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
	})
	return render(
		<QueryClientProvider client={client}>
			<DeckDetailPage deckId="deck-1" />
		</QueryClientProvider>
	)
}

beforeEach(() => {
	vi.clearAllMocks()
	mocked.getDeck.mockResolvedValue(deck)
})

describe('DeckDetailPage', () => {
	it('shows the deck header and its cards', async () => {
		mocked.listCards.mockResolvedValue([
			card(),
			card({ id: 'card-2', front: 'ねこ', back: 'cat' }),
		])

		renderPage()

		expect(await screen.findByText('Japanese Basics')).toBeInTheDocument()
		expect(await screen.findByText('こんにちは')).toBeInTheDocument()
		expect(screen.getByText('hello')).toBeInTheDocument()
		expect(screen.getByText('ねこ')).toBeInTheDocument()
	})

	it('shows an empty state inviting the first card', async () => {
		mocked.listCards.mockResolvedValue([])

		renderPage()

		expect(await screen.findByText(/no cards yet/i)).toBeInTheDocument()
	})

	it('adds a card through the dialog', async () => {
		const user = userEvent.setup()
		mocked.listCards.mockResolvedValue([])
		mocked.createCard.mockResolvedValue(card({ front: 'みず', back: 'water' }))

		renderPage()

		await user.click(await screen.findByRole('button', { name: /add card/i }))
		await user.type(screen.getByLabelText(/front/i), 'みず')
		await user.type(screen.getByLabelText(/back/i), 'water')
		await user.click(screen.getByRole('button', { name: /save card/i }))

		await waitFor(() =>
			expect(mocked.createCard).toHaveBeenCalledWith('deck-1', {
				front: 'みず',
				back: 'water',
			})
		)
	})

	it('edits a card through the dialog', async () => {
		const user = userEvent.setup()
		mocked.listCards.mockResolvedValue([card()])
		mocked.updateCard.mockResolvedValue(card({ back: 'good day' }))

		renderPage()

		await user.click(await screen.findByRole('button', { name: /edit card こんにちは/i }))
		const backInput = screen.getByLabelText(/back/i)
		await user.clear(backInput)
		await user.type(backInput, 'good day')
		await user.click(screen.getByRole('button', { name: /save card/i }))

		await waitFor(() =>
			expect(mocked.updateCard).toHaveBeenCalledWith('card-1', {
				front: 'こんにちは',
				back: 'good day',
			})
		)
	})

	it('shows an error state when cards fail to load', async () => {
		mocked.listCards.mockRejectedValue(new Error('boom'))

		renderPage()

		expect(await screen.findByText(/could not load/i)).toBeInTheDocument()
	})

	it('does not call the API when the add-card form is submitted empty', async () => {
		const user = userEvent.setup()
		mocked.listCards.mockResolvedValue([])

		renderPage()

		await user.click(await screen.findByRole('button', { name: /add card/i }))
		await user.click(screen.getByRole('button', { name: /save card/i }))

		expect(mocked.createCard).not.toHaveBeenCalled()
	})

	it('keeps the card when deletion is cancelled', async () => {
		const user = userEvent.setup()
		mocked.listCards.mockResolvedValue([card()])

		renderPage()

		await user.click(await screen.findByRole('button', { name: /delete card こんにちは/i }))
		await user.click(screen.getByRole('button', { name: /cancel/i }))

		expect(mocked.deleteCard).not.toHaveBeenCalled()
		expect(screen.getByText('こんにちは')).toBeInTheDocument()
	})

	it('links back to the decks list', async () => {
		mocked.listCards.mockResolvedValue([])

		renderPage()

		const back = await screen.findByRole('link', { name: /back to decks/i })
		expect(back).toHaveAttribute('href', '/pollyglot/decks')
	})

	it('deletes a card after confirmation', async () => {
		const user = userEvent.setup()
		mocked.listCards.mockResolvedValue([card()])
		mocked.deleteCard.mockResolvedValue(undefined)

		renderPage()

		await user.click(await screen.findByRole('button', { name: /delete card こんにちは/i }))
		await user.click(screen.getByRole('button', { name: /^delete$/i }))

		await waitFor(() => expect(mocked.deleteCard).toHaveBeenCalledWith('card-1'))
	})
})
