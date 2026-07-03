import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ConversationFactory, renderWithQuery } from '@/src/lib/test-utils'
import { ConversationPage } from './conversation-page'

const pushMock = vi.fn()

vi.mock('next/navigation', () => ({
	useRouter: () => ({ push: pushMock }),
}))

vi.mock('@/src/domain/services/conversation.service', () => ({
	conversationService: {
		listConversations: vi.fn(),
		createConversation: vi.fn(),
	},
}))

import { conversationService } from '@/src/domain/services/conversation.service'

const mocked = vi.mocked(conversationService)

const renderPage = () => renderWithQuery(<ConversationPage />)

beforeEach(() => {
	vi.clearAllMocks()
})

describe('ConversationPage (index)', () => {
	it('links each conversation to its chat', async () => {
		mocked.listConversations.mockResolvedValue([
			ConversationFactory.build({ id: 'conv-1', title: 'Practice Japanese' }),
			ConversationFactory.build({ id: 'conv-2', title: 'Ordering food', language: 'Spanish' }),
		])

		renderPage()

		const link = await screen.findByRole('link', { name: /practice japanese/i })
		expect(link).toHaveAttribute('href', '/pollyglot/conversation/conv-1')
		expect(screen.getByRole('link', { name: /ordering food/i })).toHaveAttribute(
			'href',
			'/pollyglot/conversation/conv-2'
		)
	})

	it('starts a conversation and navigates into it', async () => {
		const user = userEvent.setup()
		mocked.listConversations.mockResolvedValue([])
		mocked.createConversation.mockResolvedValue(
			ConversationFactory.build({ id: 'conv-9', language: 'Spanish' })
		)

		renderPage()

		await user.click(await screen.findByRole('button', { name: /new conversation/i }))
		await user.type(screen.getByLabelText(/language/i), 'Spanish')
		await user.click(screen.getByRole('button', { name: /start practicing/i }))

		await waitFor(() =>
			expect(mocked.createConversation).toHaveBeenCalledWith({ language: 'Spanish' })
		)
		expect(pushMock).toHaveBeenCalledWith('/pollyglot/conversation/conv-9')
	})

	it('does not create without a language', async () => {
		const user = userEvent.setup()
		mocked.listConversations.mockResolvedValue([])

		renderPage()

		await user.click(await screen.findByRole('button', { name: /new conversation/i }))
		await user.click(screen.getByRole('button', { name: /start practicing/i }))

		expect(mocked.createConversation).not.toHaveBeenCalled()
	})

	it('shows an empty state inviting the first conversation', async () => {
		mocked.listConversations.mockResolvedValue([])

		renderPage()

		expect(await screen.findByText(/no conversations yet/i)).toBeInTheDocument()
	})

	it('shows an error state when the list fails to load', async () => {
		mocked.listConversations.mockRejectedValue(new Error('boom'))

		renderPage()

		expect(await screen.findByText(/could not load/i)).toBeInTheDocument()
	})
})
