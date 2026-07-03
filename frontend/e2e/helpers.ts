import { expect, Page } from '@playwright/test'

export const seeded = {
	email: 'demo@pollyglot.dev',
	password: 'Password123!',
}

export async function login(page: Page, email = seeded.email, password = seeded.password) {
	await page.goto('/auth/login')
	await page.getByLabel('Email').fill(email)
	await page.getByLabel('Password', { exact: true }).fill(password)
	await page.getByRole('button', { name: /sign in|login/i }).click()
	await expect(page).toHaveURL(/\/home/)
}

export function uniqueName(prefix: string): string {
	return `${prefix} ${Date.now()}-${Math.floor(Math.random() * 1e4)}`
}
