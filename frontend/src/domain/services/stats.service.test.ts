import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('./api-client', () => ({
	default: {
		get: vi.fn(),
		put: vi.fn(),
	},
}))

import apiClient from './api-client'
import { statsService } from './stats.service'

const mocked = vi.mocked(apiClient)

beforeEach(() => {
	vi.clearAllMocks()
})

describe('statsService request contract', () => {
	it('fetches stats via GET /v1/stats', async () => {
		const stats = {
			reviews_today: 5,
			daily_goal: 20,
			total_reviews: 42,
			unique_cards: 17,
			streak_days: 3,
			reviews_per_day: [],
		}
		mocked.get.mockResolvedValue({ data: stats })

		const result = await statsService.getStats()

		expect(mocked.get).toHaveBeenCalledWith('/v1/stats')
		expect(result).toEqual(stats)
	})

	it('sets the goal via PUT /v1/stats/goal', async () => {
		mocked.put.mockResolvedValue({ data: { daily_goal: 35 } })

		await statsService.setGoal(35)

		expect(mocked.put).toHaveBeenCalledWith('/v1/stats/goal', { goal: 35 })
	})
})
