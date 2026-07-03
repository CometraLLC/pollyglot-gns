import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { CardFactory, DeckFactory, renderWithQuery } from '@/src/lib/test-utils'
import { StudySessionPage } from './study-session-page'

vi.mock('@/src/domain/services/decks.service', () => ({
	decksService: {
		getDeck: vi.fn(),
		getStudyQueue: vi.fn(),
		reviewCard: vi.fn(),
	},
}))

import { decksService } from '@/src/domain/services/decks.service'

const mocked = vi.mocked(decksService)

const deck = DeckFactory.build({ id: 'deck-1', card_count: 2 })

const card = CardFactory.build

const twoCards = [card(), card({ id: 'card-2', front: 'ねこ', back: 'cat' })]

const renderPage = () => renderWithQuery(<StudySessionPage deckId="deck-1" />)

beforeEach(() => {
	vi.clearAllMocks()
	mocked.getDeck.mockResolvedValue(deck)
	mocked.reviewCard.mockImplementation(async (id) => card({ id }))
})

describe('StudySessionPage', () => {
	it('shows a loading state while the queue fetches', () => {
		mocked.getStudyQueue.mockReturnValue(new Promise(() => {}))

		renderPage()

		expect(screen.getByText(/loading/i)).toBeInTheDocument()
	})

	it('shows an error state when the queue fails to load', async () => {
		mocked.getStudyQueue.mockRejectedValue(new Error('boom'))

		renderPage()

		expect(await screen.findByText(/could not load/i)).toBeInTheDocument()
	})

	it('celebrates when there is nothing due', async () => {
		mocked.getStudyQueue.mockResolvedValue([])

		renderPage()

		expect(await screen.findByText(/all caught up/i)).toBeInTheDocument()
		expect(screen.getByRole('link', { name: /back to decks/i })).toHaveAttribute(
			'href',
			'/pollyglot/decks'
		)
	})

	it('starts face down with the answer hidden (original flashcard contract)', async () => {
		mocked.getStudyQueue.mockResolvedValue(twoCards)

		renderPage()

		const flip = await screen.findByRole('button', { name: /show answer/i })
		expect(flip).toHaveAttribute('aria-pressed', 'false')
		expect(screen.getByText('こんにちは')).toBeInTheDocument()
	})

	it('keeps rating buttons out of the tab order until flipped', async () => {
		mocked.getStudyQueue.mockResolvedValue(twoCards)

		renderPage()

		const gotIt = await screen.findByRole('button', { name: 'Rate as Got it!' })
		expect(gotIt).toHaveAttribute('tabindex', '-1')

		await userEvent.setup().click(screen.getByRole('button', { name: /show answer/i }))
		expect(gotIt).toHaveAttribute('tabindex', '0')
	})

	it('flips to reveal the answer and offers all five ratings', async () => {
		const user = userEvent.setup()
		mocked.getStudyQueue.mockResolvedValue(twoCards)

		renderPage()

		await user.click(await screen.findByRole('button', { name: /show answer/i }))

		expect(screen.getByRole('button', { name: /hide answer/i })).toHaveAttribute(
			'aria-pressed',
			'true'
		)
		expect(screen.getByText('hello')).toBeInTheDocument()
		for (const label of ['Forgot', 'Difficult', 'Okay', 'Almost', 'Got it!']) {
			expect(screen.getByRole('button', { name: `Rate as ${label}` })).toBeInTheDocument()
		}
	})

	it('sends the right rating and advances to the next card face down', async () => {
		const user = userEvent.setup()
		mocked.getStudyQueue.mockResolvedValue(twoCards)

		renderPage()

		await user.click(await screen.findByRole('button', { name: /show answer/i }))
		await user.click(screen.getByRole('button', { name: 'Rate as Got it!' }))

		await waitFor(() => expect(mocked.reviewCard).toHaveBeenCalledWith('card-1', 4))
		expect(await screen.findByText('ねこ')).toBeInTheDocument()
		expect(screen.getByRole('button', { name: /show answer/i })).toHaveAttribute(
			'aria-pressed',
			'false'
		)
		expect(screen.queryByText('cat')).not.toBeInTheDocument()
	})

	it('maps every rating button to its SM-2 value', async () => {
		const user = userEvent.setup()
		const cards = ['a', 'b', 'c', 'd', 'e'].map((id, i) =>
			card({ id, front: `front-${i}`, back: `back-${i}` })
		)
		mocked.getStudyQueue.mockResolvedValue(cards)

		renderPage()

		const ratings: Array<[string, number]> = [
			['Rate as Forgot', 0],
			['Rate as Difficult', 1],
			['Rate as Okay', 2],
			['Rate as Almost', 3],
			['Rate as Got it!', 4],
		]
		for (const [label, value] of ratings) {
			await user.click(await screen.findByRole('button', { name: /show answer/i }))
			await user.click(screen.getByRole('button', { name: label }))
			await waitFor(() =>
				expect(mocked.reviewCard).toHaveBeenCalledWith(expect.any(String), value)
			)
		}
		expect(mocked.reviewCard).toHaveBeenCalledTimes(5)
	})

	it('tracks progress and cards flipped', async () => {
		const user = userEvent.setup()
		mocked.getStudyQueue.mockResolvedValue(twoCards)

		renderPage()

		expect(await screen.findByText(/1 of 2/)).toBeInTheDocument()
		expect(screen.getByLabelText(/cards flipped: 0/i)).toBeInTheDocument()

		// flip forth and back: two flips
		await user.click(screen.getByRole('button', { name: /show answer/i }))
		await user.click(screen.getByRole('button', { name: /hide answer/i }))
		expect(screen.getByLabelText(/cards flipped: 2/i)).toBeInTheDocument()

		await user.click(screen.getByRole('button', { name: /show answer/i }))
		await user.click(screen.getByRole('button', { name: 'Rate as Okay' }))
		expect(await screen.findByText(/2 of 2/)).toBeInTheDocument()
	})

	it('finishes the session with a summary after the last card', async () => {
		const user = userEvent.setup()
		mocked.getStudyQueue.mockResolvedValue([card()])

		renderPage()

		await user.click(await screen.findByRole('button', { name: /show answer/i }))
		await user.click(screen.getByRole('button', { name: 'Rate as Almost' }))

		expect(await screen.findByText(/session complete/i)).toBeInTheDocument()
		expect(screen.getByText(/reviewed 1 card/i)).toBeInTheDocument()
		expect(screen.getByRole('link', { name: /back to decks/i })).toBeInTheDocument()
	})

	it('blanks cloze deletions on the front and reveals them on flip', async () => {
		const user = userEvent.setup()
		mocked.getStudyQueue.mockResolvedValue([
			card({ card_type: 'cloze', front: '水を{{c1::飲みます}}', back: 'drink water' }),
		])

		renderPage()

		expect(await screen.findByText('水を[…]')).toBeInTheDocument()
		expect(screen.queryByText('水を飲みます')).not.toBeInTheDocument()

		await user.click(screen.getByRole('button', { name: /show answer/i }))

		expect(screen.getByText('水を飲みます')).toBeInTheDocument()
		expect(screen.getByText('drink water')).toBeInTheDocument()
	})

	it('offers pronunciation when the browser can speak', async () => {
		const speakMock = vi.fn()
		vi.stubGlobal('speechSynthesis', { speak: speakMock, cancel: vi.fn() })
		vi.stubGlobal(
			'SpeechSynthesisUtterance',
			class {
				text: string
				lang = ''
				constructor(text: string) {
					this.text = text
				}
			}
		)
		const user = userEvent.setup()
		mocked.getStudyQueue.mockResolvedValue(twoCards)

		renderPage()

		await user.click(await screen.findByRole('button', { name: /pronounce/i }))

		expect(speakMock).toHaveBeenCalledTimes(1)
		expect(speakMock.mock.calls[0][0].text).toBe('こんにちは')
		expect(speakMock.mock.calls[0][0].lang).toBe('ja-JP')
		vi.unstubAllGlobals()
	})

	it('stays on the card and surfaces an error when the review fails', async () => {
		const user = userEvent.setup()
		mocked.getStudyQueue.mockResolvedValue(twoCards)
		mocked.reviewCard.mockRejectedValue(new Error('offline'))

		renderPage()

		await user.click(await screen.findByRole('button', { name: /show answer/i }))
		await user.click(screen.getByRole('button', { name: 'Rate as Okay' }))

		expect(await screen.findByText(/could not save/i)).toBeInTheDocument()
		// still on the first card, answer still showing
		expect(screen.getByText('hello')).toBeInTheDocument()
		expect(screen.queryByText('ねこ')).not.toBeInTheDocument()
	})
})
