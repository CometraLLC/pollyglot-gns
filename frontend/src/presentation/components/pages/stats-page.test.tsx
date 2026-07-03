import { screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { renderWithQuery } from '@/src/lib/test-utils'
import { StatsPage } from './stats-page'

vi.mock('@/src/domain/services/stats.service', () => ({
	statsService: {
		getStats: vi.fn(),
		setGoal: vi.fn(),
	},
}))

// jsdom has no layout engine, so Recharts cannot render; the chart is
// stubbed here and verified visually in the browser instead.
vi.mock('recharts', () => {
	const Stub = ({ children }: { children?: React.ReactNode }) => (
		<div data-testid="chart">{children}</div>
	)
	return {
		ResponsiveContainer: Stub,
		BarChart: Stub,
		Bar: () => null,
		XAxis: () => null,
		YAxis: () => null,
		Tooltip: () => null,
	}
})

import { statsService } from '@/src/domain/services/stats.service'
import type { Stats } from '@/src/domain/services/stats.service'

const mocked = vi.mocked(statsService)

function stats(overrides: Partial<Stats> = {}): Stats {
	return {
		reviews_today: 5,
		daily_goal: 20,
		total_reviews: 42,
		unique_cards: 17,
		streak_days: 3,
		reviews_per_day: [
			{ date: '2026-07-01', count: 2 },
			{ date: '2026-07-02', count: 5 },
		],
		...overrides,
	}
}

const renderPage = () => renderWithQuery(<StatsPage />)

beforeEach(() => {
	vi.clearAllMocks()
})

describe('StatsPage', () => {
	it('shows the four headline stats', async () => {
		mocked.getStats.mockResolvedValue(stats())

		renderPage()

		expect(await screen.findByLabelText(/day streak: 3/i)).toHaveTextContent('3')
		expect(screen.getByLabelText(/reviews today: 5/i)).toHaveTextContent('5')
		expect(screen.getByLabelText(/total reviews: 42/i)).toHaveTextContent('42')
		expect(screen.getByLabelText(/unique words: 17/i)).toHaveTextContent('17')
	})

	it('shows progress toward the daily goal', async () => {
		mocked.getStats.mockResolvedValue(stats({ reviews_today: 5, daily_goal: 20 }))

		renderPage()

		expect(await screen.findByText(/5 \/ 20/)).toBeInTheDocument()
		const bar = screen.getByRole('progressbar', { name: /daily goal/i })
		expect(bar).toHaveAttribute('aria-valuenow', '5')
		expect(bar).toHaveAttribute('aria-valuemax', '20')
	})

	it('celebrates when the goal is met', async () => {
		mocked.getStats.mockResolvedValue(stats({ reviews_today: 25, daily_goal: 20 }))

		renderPage()

		expect(await screen.findByText(/goal met/i)).toBeInTheDocument()
	})

	it('edits the daily goal', async () => {
		const { default: userEvent } = await import('@testing-library/user-event')
		const user = userEvent.setup()
		mocked.getStats.mockResolvedValue(stats())
		mocked.setGoal.mockResolvedValue(undefined)

		renderPage()

		await user.click(await screen.findByRole('button', { name: /change goal/i }))
		const input = screen.getByRole('spinbutton', { name: /reviews per day/i })
		await user.clear(input)
		await user.type(input, '35')
		await user.click(screen.getByRole('button', { name: /save goal/i }))

		await vi.waitFor(() => expect(mocked.setGoal).toHaveBeenCalledWith(35))
	})

	it('offers the daily series as an accessible table', async () => {
		mocked.getStats.mockResolvedValue(stats())

		renderPage()

		const table = await screen.findByRole('table', { name: /reviews per day/i })
		expect(table).toBeInTheDocument()
		expect(screen.getByRole('cell', { name: '2026-07-02' })).toBeInTheDocument()
		expect(screen.getByRole('cell', { name: '5' })).toBeInTheDocument()
	})

	it('encourages a first review when everything is zero', async () => {
		mocked.getStats.mockResolvedValue(
			stats({
				reviews_today: 0,
				total_reviews: 0,
				unique_cards: 0,
				streak_days: 0,
				reviews_per_day: [],
			})
		)

		renderPage()

		expect(await screen.findByText(/no reviews yet/i)).toBeInTheDocument()
	})

	it('shows an error state when stats fail to load', async () => {
		mocked.getStats.mockRejectedValue(new Error('boom'))

		renderPage()

		expect(await screen.findByText(/could not load/i)).toBeInTheDocument()
	})

	it('shows a loading state first', () => {
		mocked.getStats.mockReturnValue(new Promise(() => {}))

		renderPage()

		expect(screen.getByText(/loading/i)).toBeInTheDocument()
	})
})
