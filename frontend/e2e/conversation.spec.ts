import { expect, test } from '@playwright/test'
import { login } from './helpers'

test('start a conversation and exchange with the Socratic tutor', async ({ page }) => {
	await login(page)

	await page.goto('/pollyglot/conversation')
	await page.getByRole('button', { name: /new conversation/i }).click()
	await page.getByLabel('Language').fill('Japanese')
	await page.getByRole('button', { name: /start practicing/i }).click()

	// tutor greets first, always ending with a question
	const greeting = page.getByLabel(/tutor said/i).first()
	await expect(greeting).toBeVisible()
	await expect(greeting).toContainText('?')

	// send a message; the tutor probes back, quoting the learner
	await page.getByLabel('Message', { exact: true }).fill('こんにちは')
	await page.getByRole('button', { name: 'Send', exact: true }).click()

	await expect(page.getByLabel(/you said/i)).toContainText('こんにちは')
	const reply = page.getByLabel(/tutor said/i).last()
	await expect(reply).toContainText('こんにちは')
	await expect(reply).toContainText('?')
})
