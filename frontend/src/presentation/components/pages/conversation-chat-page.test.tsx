import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import {
	ConversationFactory,
	MessageFactory,
	renderWithQuery,
} from '@/src/lib/test-utils'
import { ConversationChatPage } from './conversation-chat-page'

vi.mock('@/src/domain/services/conversation.service', () => ({
	conversationService: {
		listConversations: vi.fn(),
		listMessages: vi.fn(),
		sendMessage: vi.fn(),
	},
}))

import { conversationService } from '@/src/domain/services/conversation.service'

const mocked = vi.mocked(conversationService)

const conversation = ConversationFactory.build({ id: 'conv-1', title: 'Practice Japanese' })

const renderPage = () => renderWithQuery(<ConversationChatPage conversationId="conv-1" />)

beforeEach(() => {
	vi.clearAllMocks()
	mocked.listConversations.mockResolvedValue([conversation])
})

describe('ConversationChatPage', () => {
	it('renders the history with tutor and learner turns distinguishable', async () => {
		mocked.listMessages.mockResolvedValue([
			MessageFactory.build({ id: 'm1', role: 'tutor', content: 'What do you know?' }),
			MessageFactory.build({ id: 'm2', role: 'user', content: 'こんにちは' }),
		])

		renderPage()

		expect(await screen.findByText('What do you know?')).toBeInTheDocument()
		expect(screen.getByText('こんにちは')).toBeInTheDocument()
		expect(screen.getByLabelText(/tutor said/i)).toHaveTextContent('What do you know?')
		expect(screen.getByLabelText(/you said/i)).toHaveTextContent('こんにちは')
	})

	it('shows an error state when history fails to load', async () => {
		mocked.listMessages.mockRejectedValue(new Error('boom'))

		renderPage()

		expect(await screen.findByText(/could not load/i)).toBeInTheDocument()
	})

	it('sends a message and appends the whole exchange without refetching', async () => {
		const user = userEvent.setup()
		mocked.listMessages.mockResolvedValue([
			MessageFactory.build({ id: 'm1', role: 'tutor', content: 'Welcome?' }),
		])
		mocked.sendMessage.mockResolvedValue({
			user_message: MessageFactory.build({ id: 'm2', role: 'user', content: 'ねこ' }),
			tutor_message: MessageFactory.build({
				id: 'm3',
				role: 'tutor',
				content: 'What does ねこ mean?',
			}),
		})

		renderPage()
		await screen.findByText('Welcome?')

		await user.type(screen.getByLabelText(/message/i), 'ねこ')
		await user.click(screen.getByRole('button', { name: /send/i }))

		expect(await screen.findByText('ねこ')).toBeInTheDocument()
		expect(screen.getByText('What does ねこ mean?')).toBeInTheDocument()
		expect(mocked.sendMessage).toHaveBeenCalledWith('conv-1', 'ねこ')
		expect(mocked.listMessages).toHaveBeenCalledTimes(1)

		// input clears after a successful send
		expect(screen.getByLabelText(/message/i)).toHaveValue('')
	})

	it('does not send empty messages', async () => {
		const user = userEvent.setup()
		mocked.listMessages.mockResolvedValue([])

		renderPage()
		await screen.findByLabelText(/message/i)

		await user.click(screen.getByRole('button', { name: /send/i }))

		expect(mocked.sendMessage).not.toHaveBeenCalled()
	})

	it('keeps the draft and shows an error when sending fails', async () => {
		const user = userEvent.setup()
		mocked.listMessages.mockResolvedValue([])
		mocked.sendMessage.mockRejectedValue(new Error('offline'))

		renderPage()
		await screen.findByLabelText(/message/i)

		await user.type(screen.getByLabelText(/message/i), 'ねこ')
		await user.click(screen.getByRole('button', { name: /send/i }))

		expect(await screen.findByText(/could not send/i)).toBeInTheDocument()
		expect(screen.getByLabelText(/message/i)).toHaveValue('ねこ')
	})

	it('links back to the conversation list', async () => {
		mocked.listMessages.mockResolvedValue([])

		renderPage()

		const back = await screen.findByRole('link', { name: /conversations/i })
		expect(back).toHaveAttribute('href', '/pollyglot/conversation')
	})
})
