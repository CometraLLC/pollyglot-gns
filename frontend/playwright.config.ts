import { defineConfig, devices } from '@playwright/test'

// E2E suite drives the real stack: backend API on :8080 (Docker locally,
// service containers + `go run` in CI) and the Next app on :3000.
// Selectors are role/label only — never CSS classes (Pollyglot#26).
export default defineConfig({
	testDir: './e2e',
	fullyParallel: false,
	retries: process.env.CI ? 1 : 0,
	workers: 1,
	reporter: process.env.CI ? [['github'], ['list']] : 'list',
	timeout: 30_000,
	use: {
		baseURL: 'http://localhost:3000',
		trace: 'on-first-retry',
	},
	projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
	webServer: {
		command: 'bun run dev',
		url: 'http://localhost:3000',
		reuseExistingServer: !process.env.CI,
		timeout: 120_000,
	},
})
