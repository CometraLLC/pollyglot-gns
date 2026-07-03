import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('./api-client', () => ({
	default: {
		post: vi.fn(),
	},
}))

import apiClient from './api-client'
import { translateService } from './translate.service'

const mocked = vi.mocked(apiClient)

beforeEach(() => {
	vi.clearAllMocks()
})

describe('translateService request contract', () => {
	it('translates via POST /v1/translate with the exact payload', async () => {
		mocked.post.mockResolvedValue({
			data: { text: 'hello', from: 'English', to: 'Spanish', translation: 'hola' },
		})

		const result = await translateService.translate({
			text: 'hello',
			from: 'English',
			to: 'Spanish',
		})

		expect(mocked.post).toHaveBeenCalledWith('/v1/translate', {
			text: 'hello',
			from: 'English',
			to: 'Spanish',
		})
		expect(result.translation).toBe('hola')
	})
})
