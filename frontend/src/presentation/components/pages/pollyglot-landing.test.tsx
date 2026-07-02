import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'
import { PollyglotLanding } from './pollyglot-landing'

describe('PollyglotLanding', () => {
	it('renders the Pollyglot brand and auth calls to action', () => {
		render(<PollyglotLanding />)

		expect(screen.getAllByText('Pollyglot').length).toBeGreaterThan(0)
		expect(screen.getAllByRole('link', { name: /start learning/i }).length).toBeGreaterThan(0)
		expect(screen.getAllByRole('link', { name: /sign in/i }).length).toBeGreaterThan(0)
	})

	describe('hero flashcard', () => {
		it('starts face down with the answer hidden', () => {
			render(<PollyglotLanding />)

			const card = screen.getByRole('button', { name: /show answer/i })
			expect(card).toHaveAttribute('aria-pressed', 'false')
		})

		it('flips to reveal the answer when clicked', async () => {
			const user = userEvent.setup()
			render(<PollyglotLanding />)

			await user.click(screen.getByRole('button', { name: /show answer/i }))

			const flipped = screen.getByRole('button', { name: /hide answer/i })
			expect(flipped).toHaveAttribute('aria-pressed', 'true')
		})

		it('keeps rating buttons out of the tab order until the answer is showing', async () => {
			const user = userEvent.setup()
			render(<PollyglotLanding />)

			const rating = screen.getByRole('button', { name: 'Rate as Got it!' })
			expect(rating).toHaveAttribute('tabindex', '-1')

			await user.click(screen.getByRole('button', { name: /show answer/i }))
			expect(rating).toHaveAttribute('tabindex', '0')
		})

		it('offers all five ratings from Forgot to Got it!', async () => {
			const user = userEvent.setup()
			render(<PollyglotLanding />)

			await user.click(screen.getByRole('button', { name: /show answer/i }))

			for (const label of ['Forgot', 'Difficult', 'Okay', 'Almost', 'Got it!']) {
				expect(screen.getByRole('button', { name: `Rate as ${label}` })).toBeInTheDocument()
			}
		})

		it('advances to the next card after rating', async () => {
			const user = userEvent.setup()
			render(<PollyglotLanding />)

			const firstWord = 'こんにちは'
			expect(screen.getByText(firstWord)).toBeInTheDocument()

			await user.click(screen.getByRole('button', { name: /show answer/i }))
			await user.click(screen.getByRole('button', { name: 'Rate as Okay' }))

			// The card flips back and swaps content after a short delay.
			expect(await screen.findByText('gato')).toBeInTheDocument()
			expect(screen.getByRole('button', { name: /show answer/i })).toHaveAttribute(
				'aria-pressed',
				'false'
			)
		})
	})
})
