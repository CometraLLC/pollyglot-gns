import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { Card } from '@/src/domain/services/decks.service'
import { CardFactory, DeckFactory, renderWithQuery } from '@/src/lib/test-utils'
import { DeckDetailPage } from './deck-detail-page'

vi.mock('@/src/domain/services/decks.service', () => ({
	decksService: {
		getDeck: vi.fn(),
		listCards: vi.fn(),
		createCard: vi.fn(),
		updateCard: vi.fn(),
		deleteCard: vi.fn(),
		exportDeck: vi.fn(),
		importDeck: vi.fn(),
		shareDeck: vi.fn(),
		unshareDeck: vi.fn(),
	},
}))

import { decksService } from '@/src/domain/services/decks.service'

const mocked = vi.mocked(decksService)

const deck = DeckFactory.build({ id: 'deck-1', card_count: 2 })

const card = (overrides: Partial<Card> = {}): Card =>
	CardFactory.build({ id: 'card-1', deck_id: 'deck-1', ...overrides })

const renderPage = () => renderWithQuery(<DeckDetailPage deckId="deck-1" />)

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
		await user.type(screen.getByLabelText(/^front$/i), 'みず')
		await user.type(screen.getByLabelText(/^back$/i), 'water')
		await user.click(screen.getByRole('button', { name: /save card/i }))

		await waitFor(() =>
			expect(mocked.createCard).toHaveBeenCalledWith('deck-1', {
				front: 'みず',
				back: 'water',
			})
		)
	})

	it('creates a cloze card through the type selector', async () => {
		const user = userEvent.setup()
		mocked.listCards.mockResolvedValue([])
		mocked.createCard.mockResolvedValue(card({ card_type: 'cloze' }))

		renderPage()

		await user.click(await screen.findByRole('button', { name: /add card/i }))
		await user.selectOptions(screen.getByLabelText(/card type/i), 'cloze')
		await user.click(screen.getByLabelText(/^front$/i))
		await user.paste('水を{{c1::飲みます}}')
		await user.type(screen.getByLabelText(/^back$/i), 'drink water')
		await user.click(screen.getByRole('button', { name: /save card/i }))

		await waitFor(() =>
			expect(mocked.createCard).toHaveBeenCalledWith('deck-1', {
				front: '水を{{c1::飲みます}}',
				back: 'drink water',
				card_type: 'cloze',
			})
		)
	})

	it('reverse checkbox requests the mirrored card', async () => {
		const user = userEvent.setup()
		mocked.listCards.mockResolvedValue([])
		mocked.createCard.mockResolvedValue(card())

		renderPage()

		await user.click(await screen.findByRole('button', { name: /add card/i }))
		await user.type(screen.getByLabelText(/^front$/i), 'ねこ')
		await user.type(screen.getByLabelText(/^back$/i), 'cat')
		await user.click(screen.getByLabelText(/also create reversed card/i))
		await user.click(screen.getByRole('button', { name: /save card/i }))

		await waitFor(() =>
			expect(mocked.createCard).toHaveBeenCalledWith('deck-1', {
				front: 'ねこ',
				back: 'cat',
				reverse: true,
			})
		)
	})

	it('marks cloze cards in the list', async () => {
		mocked.listCards.mockResolvedValue([card({ card_type: 'cloze', front: '水を{{c1::飲みます}}' })])

		renderPage()

		expect(await screen.findByText('cloze')).toBeInTheDocument()
	})

	it('edits a card through the dialog', async () => {
		const user = userEvent.setup()
		mocked.listCards.mockResolvedValue([card()])
		mocked.updateCard.mockResolvedValue(card({ back: 'good day' }))

		renderPage()

		await user.click(await screen.findByRole('button', { name: /edit card こんにちは/i }))
		const backInput = screen.getByLabelText(/^back$/i)
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

	it('imports a file and shows the summary with row errors', async () => {
		const user = userEvent.setup()
		mocked.listCards.mockResolvedValue([])
		mocked.importDeck.mockResolvedValue({
			imported: 2,
			skipped: [{ line: 3, error: 'front is required' }],
		})

		renderPage()

		await user.click(await screen.findByRole('button', { name: /^import$/i }))
		const file = new File(['front,back\nねこ,cat'], 'cards.csv', { type: 'text/csv' })
		await user.upload(screen.getByLabelText(/csv or tsv file/i), file)
		await user.click(screen.getByRole('button', { name: /import cards/i }))

		expect(await screen.findByText(/imported 2 cards/i)).toBeInTheDocument()
		expect(screen.getByText(/line 3: front is required/i)).toBeInTheDocument()
		expect(mocked.importDeck).toHaveBeenCalledWith('deck-1', file)
	})

	it('does not import without a file', async () => {
		const user = userEvent.setup()
		mocked.listCards.mockResolvedValue([])

		renderPage()

		await user.click(await screen.findByRole('button', { name: /^import$/i }))
		await user.click(screen.getByRole('button', { name: /import cards/i }))

		expect(mocked.importDeck).not.toHaveBeenCalled()
	})

	it('exports the deck as csv', async () => {
		const user = userEvent.setup()
		mocked.listCards.mockResolvedValue([])
		mocked.exportDeck.mockResolvedValue(new Blob(['front,back,card_type']))
		vi.stubGlobal('URL', {
			createObjectURL: vi.fn(() => 'blob:x'),
			revokeObjectURL: vi.fn(),
		})

		renderPage()

		await user.click(await screen.findByRole('button', { name: /export csv/i }))

		await waitFor(() => expect(mocked.exportDeck).toHaveBeenCalledWith('deck-1', 'csv'))
		vi.unstubAllGlobals()
	})

	it('enables sharing and shows the link', async () => {
		const user = userEvent.setup()
		mocked.listCards.mockResolvedValue([])
		mocked.shareDeck.mockResolvedValue({ share_code: 'ABCDEF2345' })

		renderPage()

		await user.click(await screen.findByRole('button', { name: /^share$/i }))
		await user.click(screen.getByRole('button', { name: /enable sharing/i }))

		await waitFor(() => expect(mocked.shareDeck).toHaveBeenCalledWith('deck-1'))
		expect(await screen.findByText(/\/pollyglot\/shared\/ABCDEF2345/)).toBeInTheDocument()
	})

	it('disables sharing', async () => {
		const user = userEvent.setup()
		mocked.getDeck.mockResolvedValue(DeckFactory.build({ id: 'deck-1', card_count: 2, share_code: 'ABCDEF2345' }))
		mocked.listCards.mockResolvedValue([])
		mocked.unshareDeck.mockResolvedValue(undefined)

		renderPage()

		await user.click(await screen.findByRole('button', { name: /^share$/i }))
		await user.click(screen.getByRole('button', { name: /disable sharing/i }))

		await waitFor(() => expect(mocked.unshareDeck).toHaveBeenCalledWith('deck-1'))
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
