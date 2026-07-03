import { NextIntlClientProvider } from 'next-intl'
import { render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import en from '@/locales/en.json'
import { Sidebar } from './sidebar'

vi.mock('next/navigation', () => ({
	usePathname: () => '/pollyglot',
}))

function renderSidebar() {
	return render(
		<NextIntlClientProvider locale="en" messages={en}>
			<Sidebar isOpen onClose={() => {}} />
		</NextIntlClientProvider>
	)
}

describe('Sidebar navigation', () => {
	it('links every Pollyglot product page', () => {
		renderSidebar()

		const expected: Array<[string, string]> = [
			['Decks', '/pollyglot/decks'],
			['Study', '/pollyglot/study'],
			['Translate', '/pollyglot/translate'],
			['Conversation', '/pollyglot/conversation'],
			['Progress', '/pollyglot/stats'],
		]
		for (const [label, href] of expected) {
			const link = screen.getByRole('link', { name: label })
			expect(link).toHaveAttribute('href', href)
		}
	})

	it('keeps the general navigation', () => {
		renderSidebar()

		expect(screen.getByRole('link', { name: 'Home' })).toHaveAttribute('href', '/home')
		expect(screen.getByRole('link', { name: 'Documentation' })).toHaveAttribute('href', '/docs')
	})

	it('brands as Pollyglot in the header and links the hub', () => {
		renderSidebar()

		// header brand + the hub nav item
		expect(screen.getAllByText('Pollyglot').length).toBeGreaterThanOrEqual(2)
		expect(screen.getByRole('link', { name: 'Pollyglot' })).toHaveAttribute('href', '/pollyglot')
	})
})
