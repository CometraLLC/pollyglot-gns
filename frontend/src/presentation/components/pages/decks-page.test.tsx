import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { Deck } from '@/src/domain/services/decks.service'
import { DecksPage } from './decks-page'

vi.mock('@/src/domain/services/decks.service', () => ({
	decksService: {
		listDecks: vi.fn(),
		createDeck: vi.fn(),
		updateDeck: vi.fn(),
		deleteDeck: vi.fn(),
	},
}))

import { decksService } from '@/src/domain/services/decks.service'

const mocked = vi.mocked(decksService)

function deck(overrides: Partial<Deck> = {}): Deck {
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

function renderPage() {
	const client = new QueryClient({
		defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
	})
	return render(
		<QueryClientProvider client={client}>
			<DecksPage />
		</QueryClientProvider>
	)
}

beforeEach(() => {
	vi.clearAllMocks()
})

describe('DecksPage', () => {
	it('lists decks with language pair and card count', async () => {
		mocked.listDecks.mockResolvedValue([
			deck(),
			deck({ id: 'deck-2', name: 'Spanish Verbs', source_lang: 'Spanish', card_count: 1 }),
		])

		renderPage()

		expect(await screen.findByText('Japanese Basics')).toBeInTheDocument()
		expect(screen.getByText('Spanish Verbs')).toBeInTheDocument()
		expect(screen.getByText(/Japanese → English/)).toBeInTheDocument()
		expect(screen.getByText(/3 cards/)).toBeInTheDocument()
		expect(screen.getByText(/1 card\b/)).toBeInTheDocument()
	})

	it('shows an empty state inviting the first deck', async () => {
		mocked.listDecks.mockResolvedValue([])

		renderPage()

		expect(await screen.findByText(/no decks yet/i)).toBeInTheDocument()
	})

	it('creates a deck through the dialog', async () => {
		const user = userEvent.setup()
		mocked.listDecks.mockResolvedValue([])
		mocked.createDeck.mockResolvedValue(deck({ name: 'French Food' }))

		renderPage()

		await user.click(await screen.findByRole('button', { name: /new deck/i }))
		await user.type(screen.getByLabelText(/^name$/i), 'French Food')
		await user.type(screen.getByLabelText(/learning language/i), 'French')
		await user.type(screen.getByLabelText(/your language/i), 'English')
		await user.click(screen.getByRole('button', { name: /create deck/i }))

		await waitFor(() =>
			expect(mocked.createDeck).toHaveBeenCalledWith({
				name: 'French Food',
				source_lang: 'French',
				target_lang: 'English',
			})
		)
	})

	it('deletes a deck after confirmation', async () => {
		const user = userEvent.setup()
		mocked.listDecks.mockResolvedValue([deck()])
		mocked.deleteDeck.mockResolvedValue(undefined)

		renderPage()

		await user.click(await screen.findByRole('button', { name: /delete japanese basics/i }))
		// confirmation dialog
		await user.click(screen.getByRole('button', { name: /^delete$/i }))

		await waitFor(() => expect(mocked.deleteDeck).toHaveBeenCalledWith('deck-1'))
	})

	it('shows a loading indicator while decks are fetching', async () => {
		let resolve!: (value: Deck[]) => void
		mocked.listDecks.mockReturnValue(new Promise<Deck[]>((r) => (resolve = r)))

		renderPage()

		expect(screen.getByText(/loading/i)).toBeInTheDocument()
		resolve([])
		expect(await screen.findByText(/no decks yet/i)).toBeInTheDocument()
	})

	it('shows an error state with the failure surfaced', async () => {
		mocked.listDecks.mockRejectedValue(new Error('network down'))

		renderPage()

		expect(await screen.findByText(/could not load your decks/i)).toBeInTheDocument()
	})

	it('does not call the API when the create form is submitted empty', async () => {
		const user = userEvent.setup()
		mocked.listDecks.mockResolvedValue([])

		renderPage()

		await user.click(await screen.findByRole('button', { name: /new deck/i }))
		await user.click(screen.getByRole('button', { name: /create deck/i }))

		expect(mocked.createDeck).not.toHaveBeenCalled()
	})

	it('keeps the deck when deletion is cancelled', async () => {
		const user = userEvent.setup()
		mocked.listDecks.mockResolvedValue([deck()])

		renderPage()

		await user.click(await screen.findByRole('button', { name: /delete japanese basics/i }))
		await user.click(screen.getByRole('button', { name: /cancel/i }))

		expect(mocked.deleteDeck).not.toHaveBeenCalled()
		expect(screen.getByText('Japanese Basics')).toBeInTheDocument()
	})

	it('links each deck to its detail page', async () => {
		mocked.listDecks.mockResolvedValue([deck()])

		renderPage()

		const link = await screen.findByRole('link', { name: /japanese basics/i })
		expect(link).toHaveAttribute('href', '/pollyglot/decks/deck-1')
	})

	it('renames a deck through the edit dialog', async () => {
		const user = userEvent.setup()
		mocked.listDecks.mockResolvedValue([deck()])
		mocked.updateDeck.mockResolvedValue(deck({ name: 'JLPT N5' }))

		renderPage()

		await user.click(await screen.findByRole('button', { name: /edit japanese basics/i }))
		const nameInput = screen.getByLabelText(/^name$/i)
		await user.clear(nameInput)
		await user.type(nameInput, 'JLPT N5')
		await user.click(screen.getByRole('button', { name: /save/i }))

		await waitFor(() =>
			expect(mocked.updateDeck).toHaveBeenCalledWith('deck-1', {
				name: 'JLPT N5',
				source_lang: 'Japanese',
				target_lang: 'English',
			})
		)
	})
})
