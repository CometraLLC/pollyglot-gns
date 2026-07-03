import { expect, test } from '@playwright/test'
import { login } from './helpers'

test('a new account can register and land in the app', async ({ page }) => {
	const email = `e2e-${Date.now()}@pollyglot.dev`

	await page.goto('/auth/register')
	await page.getByLabel('Full Name').fill('E2E User')
	await page.getByLabel('Email').fill(email)
	await page.getByLabel('Password', { exact: true }).fill('Password123!')
	await page.getByLabel('Confirm Password').fill('Password123!')
	await page.getByRole('button', { name: 'Create Account' }).click()

	await expect(page).toHaveURL(/\/home/)
})

test('the seeded demo account signs in', async ({ page }) => {
	await login(page)

	await page.goto('/pollyglot')
	await expect(page.getByRole('heading', { name: 'Pollyglot' })).toBeVisible()
})
