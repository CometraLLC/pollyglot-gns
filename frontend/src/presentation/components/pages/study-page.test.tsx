import { screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { mockDeck as deck, renderWithQuery } from '@/src/lib/test-utils'
import { StudyPage } from './study-page'

vi.mock('@/src/domain/services/decks.service', () => ({
	decksService: {
		listDecks: vi.fn(),
	},
}))

import { decksService } from '@/src/domain/services/decks.service'

const mocked = vi.mocked(decksService)

const renderPage = () => renderWithQuery(<StudyPage />)

beforeEach(() => {
	vi.clearAllMocks()
})

describe('StudyPage (deck picker)', () => {
	it('links each deck to its study session', async () => {
		mocked.listDecks.mockResolvedValue([deck(), deck({ id: 'deck-2', name: 'Spanish Verbs' })])

		renderPage()

		const link = await screen.findByRole('link', { name: /japanese basics/i })
		expect(link).toHaveAttribute('href', '/pollyglot/study/deck-1')
		expect(screen.getByRole('link', { name: /spanish verbs/i })).toHaveAttribute(
			'href',
			'/pollyglot/study/deck-2'
		)
	})

	it('points at deck creation when there is nothing to study', async () => {
		mocked.listDecks.mockResolvedValue([])

		renderPage()

		expect(await screen.findByText(/no decks to study/i)).toBeInTheDocument()
		expect(screen.getByRole('link', { name: /create a deck/i })).toHaveAttribute(
			'href',
			'/pollyglot/decks'
		)
	})

	it('shows an error state when decks fail to load', async () => {
		mocked.listDecks.mockRejectedValue(new Error('boom'))

		renderPage()

		expect(await screen.findByText(/could not load/i)).toBeInTheDocument()
	})
})
