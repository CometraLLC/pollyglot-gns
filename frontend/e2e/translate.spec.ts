import { expect, test } from '@playwright/test'
import { login } from './helpers'

test('translate a seeded word and see the result', async ({ page }) => {
	await login(page)

	await page.goto('/pollyglot/translate')
	await page.getByLabel(/text to translate/i).fill('こんにちは')
	await page.getByRole('button', { name: 'Translate', exact: true }).click()

	await expect(page.getByText('hello', { exact: true })).toBeVisible()
	// a translation can be saved into a deck
	await expect(page.getByLabel(/save to deck/i)).toBeVisible()
})

test('unknown words get the friendly no-translation message', async ({ page }) => {
	await login(page)

	await page.goto('/pollyglot/translate')
	await page.getByLabel(/text to translate/i).fill('flibbertigibbet')
	await page.getByRole('button', { name: 'Translate', exact: true }).click()

	await expect(page.getByText(/no translation available/i)).toBeVisible()
})
