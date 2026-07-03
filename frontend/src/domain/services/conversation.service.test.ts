import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('./api-client', () => ({
	default: {
		get: vi.fn(),
		post: vi.fn(),
	},
}))

import apiClient from './api-client'
import { conversationService } from './conversation.service'

const mocked = vi.mocked(apiClient)

beforeEach(() => {
	vi.clearAllMocks()
	mocked.get.mockResolvedValue({ data: [] })
	mocked.post.mockResolvedValue({ data: {} })
})

describe('conversationService request contract', () => {
	it('lists conversations via GET /v1/conversations', async () => {
		await conversationService.listConversations()
		expect(mocked.get).toHaveBeenCalledWith('/v1/conversations')
	})

	it('creates a conversation via POST /v1/conversations', async () => {
		await conversationService.createConversation({ language: 'Japanese' })
		expect(mocked.post).toHaveBeenCalledWith('/v1/conversations', { language: 'Japanese' })
	})

	it('lists messages via GET /v1/conversations/{id}/messages', async () => {
		await conversationService.listMessages('conv-1')
		expect(mocked.get).toHaveBeenCalledWith('/v1/conversations/conv-1/messages')
	})

	it('sends a message via POST /v1/conversations/{id}/messages', async () => {
		await conversationService.sendMessage('conv-1', 'こんにちは')
		expect(mocked.post).toHaveBeenCalledWith('/v1/conversations/conv-1/messages', {
			content: 'こんにちは',
		})
	})
})
