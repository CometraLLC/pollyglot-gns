import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('./api-client', () => ({
	default: {
		post: vi.fn(),
	},
}))

import apiClient from './api-client'
import { speechService } from './speech.service'

const mocked = vi.mocked(apiClient)

beforeEach(() => {
	vi.clearAllMocks()
})

describe('speechService request contract', () => {
	it('synthesizes via POST /v1/speech expecting a blob', async () => {
		const blob = new Blob(['mp3'])
		mocked.post.mockResolvedValue({ data: blob })

		const result = await speechService.synthesize('こんにちは', 'Japanese')

		expect(mocked.post).toHaveBeenCalledWith(
			'/v1/speech',
			{ text: 'こんにちは', language: 'Japanese' },
			{ responseType: 'blob' }
		)
		expect(result).toBe(blob)
	})
})
