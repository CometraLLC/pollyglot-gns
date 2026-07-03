import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { DeckFactory, renderWithQuery } from '@/src/lib/test-utils'
import { SharedDeckPage } from './shared-deck-page'

const pushMock = vi.fn()

vi.mock('next/navigation', () => ({
	useRouter: () => ({ push: pushMock }),
}))

vi.mock('@/src/domain/services/decks.service', () => ({
	decksService: {
		getSharedDeck: vi.fn(),
		cloneSharedDeck: vi.fn(),
	},
}))

import { decksService } from '@/src/domain/services/decks.service'

const mocked = vi.mocked(decksService)

const preview = {
	name: 'Japanese Basics',
	source_lang: 'Japanese',
	target_lang: 'English',
	card_count: 6,
	sample_cards: [
		{ front: 'こんにちは', back: 'hello' },
		{ front: 'ねこ', back: 'cat' },
	],
}

const renderPage = () => renderWithQuery(<SharedDeckPage code="ABCDEF2345" />)

beforeEach(() => {
	vi.clearAllMocks()
})

describe('SharedDeckPage', () => {
	it('previews the shared deck with sample cards', async () => {
		mocked.getSharedDeck.mockResolvedValue(preview)

		renderPage()

		expect(await screen.findByText('Japanese Basics')).toBeInTheDocument()
		expect(screen.getByText(/Japanese → English/)).toBeInTheDocument()
		expect(screen.getByText(/6 cards/)).toBeInTheDocument()
		expect(screen.getByText('こんにちは')).toBeInTheDocument()
	})

	it('clones the deck and navigates to it', async () => {
		const user = userEvent.setup()
		mocked.getSharedDeck.mockResolvedValue(preview)
		mocked.cloneSharedDeck.mockResolvedValue(DeckFactory.build({ id: 'deck-9' }))

		renderPage()

		await user.click(await screen.findByRole('button', { name: /add to my decks/i }))

		await waitFor(() => expect(mocked.cloneSharedDeck).toHaveBeenCalledWith('ABCDEF2345'))
		expect(pushMock).toHaveBeenCalledWith('/pollyglot/decks/deck-9')
	})

	it('handles an unknown or disabled code', async () => {
		mocked.getSharedDeck.mockRejectedValue(new Error('404'))

		renderPage()

		expect(await screen.findByText(/no longer shared|not found/i)).toBeInTheDocument()
	})
})
