import { expect, test } from '@playwright/test'
import { login, uniqueName } from './helpers'

test('deck → card → study → review shows up in stats', async ({ page }) => {
	await login(page)
	const deckName = uniqueName('E2E Deck')

	// create a deck
	await page.goto('/pollyglot/decks')
	await page.getByRole('button', { name: /new deck/i }).click()
	await page.getByLabel('Name', { exact: true }).fill(deckName)
	await page.getByLabel(/learning language/i).fill('Japanese')
	await page.getByLabel(/your language/i).fill('English')
	await page.getByRole('button', { name: /create deck/i }).click()
	await expect(page.getByRole('link', { name: deckName })).toBeVisible()

	// add a card
	await page.getByRole('link', { name: deckName }).click()
	await page.getByRole('button', { name: /add card/i }).click()
	await page.getByLabel('Front', { exact: true }).fill('やま')
	await page.getByLabel('Back', { exact: true }).fill('mountain')
	await page.getByRole('button', { name: /save card/i }).click()
	await expect(page.getByText('やま')).toBeVisible()

	// study it: flip, then rate
	await page.goto('/pollyglot/study')
	await page.getByRole('link', { name: new RegExp(deckName) }).click()
	await page.getByRole('button', { name: /show answer/i }).click()
	await expect(page.getByText('mountain')).toBeVisible()
	await page.getByRole('button', { name: 'Rate as Got it!' }).click()
	await expect(page.getByText(/session complete/i)).toBeVisible()

	// stats reflect at least one review today
	await page.goto('/pollyglot/stats')
	const reviewsToday = page.getByLabel(/reviews today/i)
	await expect(reviewsToday).toBeVisible()
	await expect(reviewsToday).not.toHaveText(/Reviews today:?\s*0$/)
})
