import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { AxiosError, AxiosHeaders } from 'axios'
import { DeckFactory, renderWithQuery } from '@/src/lib/test-utils'
import { TranslatePage } from './translate-page'

vi.mock('@/src/domain/services/translate.service', () => ({
	translateService: {
		translate: vi.fn(),
	},
}))

vi.mock('@/src/domain/services/decks.service', () => ({
	decksService: {
		listDecks: vi.fn(),
		createCard: vi.fn(),
	},
}))

import { translateService } from '@/src/domain/services/translate.service'
import { decksService } from '@/src/domain/services/decks.service'

const mockedTranslate = vi.mocked(translateService)
const mockedDecks = vi.mocked(decksService)

const renderPage = () => renderWithQuery(<TranslatePage />)

function axios422(message: string) {
	const headers = new AxiosHeaders()
	return new AxiosError(message, '422', undefined, undefined, {
		status: 422,
		statusText: 'Unprocessable Entity',
		headers,
		config: { headers },
		data: { error: message },
	})
}

beforeEach(() => {
	vi.clearAllMocks()
	mockedDecks.listDecks.mockResolvedValue([DeckFactory.build({ id: 'deck-1', name: 'Japanese Basics' })])
})

describe('TranslatePage', () => {
	it('translates text and shows the result', async () => {
		const user = userEvent.setup()
		mockedTranslate.translate.mockResolvedValue({
			text: 'こんにちは',
			from: 'Japanese',
			to: 'English',
			translation: 'hello',
		})

		renderPage()

		await user.type(screen.getByLabelText(/text to translate/i), 'こんにちは')
		await user.clear(screen.getByLabelText(/^from$/i))
		await user.type(screen.getByLabelText(/^from$/i), 'Japanese')
		await user.clear(screen.getByLabelText(/^to$/i))
		await user.type(screen.getByLabelText(/^to$/i), 'English')
		await user.click(screen.getByRole('button', { name: /^translate$/i }))

		expect(await screen.findByText('hello')).toBeInTheDocument()
		expect(mockedTranslate.translate).toHaveBeenCalledWith({
			text: 'こんにちは',
			from: 'Japanese',
			to: 'English',
		})
	})

	it('does not call the API when the text is empty', async () => {
		const user = userEvent.setup()
		renderPage()

		await user.click(screen.getByRole('button', { name: /^translate$/i }))

		expect(mockedTranslate.translate).not.toHaveBeenCalled()
	})

	it('swaps the languages', async () => {
		const user = userEvent.setup()
		renderPage()

		const from = screen.getByLabelText(/^from$/i)
		const to = screen.getByLabelText(/^to$/i)
		await user.clear(from)
		await user.type(from, 'Japanese')
		await user.clear(to)
		await user.type(to, 'English')

		await user.click(screen.getByRole('button', { name: /swap languages/i }))

		expect(from).toHaveValue('English')
		expect(to).toHaveValue('Japanese')
	})

	it('surfaces a friendly message when no translation exists', async () => {
		const user = userEvent.setup()
		mockedTranslate.translate.mockRejectedValue(axios422('no translation available'))

		renderPage()

		await user.type(screen.getByLabelText(/text to translate/i), 'flibbertigibbet')
		await user.click(screen.getByRole('button', { name: /^translate$/i }))

		expect(await screen.findByText(/no translation available/i)).toBeInTheDocument()
	})

	it('saves the translation into a deck as a card', async () => {
		const user = userEvent.setup()
		mockedTranslate.translate.mockResolvedValue({
			text: 'こんにちは',
			from: 'Japanese',
			to: 'English',
			translation: 'hello',
		})
		mockedDecks.createCard.mockResolvedValue({
			id: 'card-9',
			deck_id: 'deck-1',
			front: 'こんにちは',
			back: 'hello',
			ease_factor: 2.5,
			interval_days: 0,
			repetitions: 0,
			due_at: '2026-07-01T00:00:00Z',
			created_at: '2026-07-01T00:00:00Z',
			updated_at: '2026-07-01T00:00:00Z',
		})

		renderPage()

		await user.type(screen.getByLabelText(/text to translate/i), 'こんにちは')
		await user.click(screen.getByRole('button', { name: /^translate$/i }))
		await screen.findByText('hello')

		await user.selectOptions(screen.getByLabelText(/save to deck/i), 'deck-1')
		await user.click(screen.getByRole('button', { name: /save as card/i }))

		await waitFor(() =>
			expect(mockedDecks.createCard).toHaveBeenCalledWith('deck-1', {
				front: 'こんにちは',
				back: 'hello',
			})
		)
		expect(await screen.findByText(/saved to japanese basics/i)).toBeInTheDocument()
	})

	it('hides the save control until there is a translation', async () => {
		renderPage()

		expect(screen.queryByRole('button', { name: /save as card/i })).not.toBeInTheDocument()
	})
})
