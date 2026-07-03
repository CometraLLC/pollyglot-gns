// Visual tour: drives the real app and saves screenshots of every page.
// Run: bun e2e/tour.ts [output-dir]   (dev stack must be running)
import { chromium, Page } from '@playwright/test'

const OUT = process.argv[2] ?? './screens'
const BASE = 'http://localhost:3000'

async function shot(page: Page, name: string) {
	await page.waitForTimeout(400) // settle animations
	await page.screenshot({ path: `${OUT}/${name}.png`, fullPage: false })
	console.log(`captured ${name}`)
}

async function main() {
	const browser = await chromium.launch()
	const page = await browser.newPage({ viewport: { width: 1360, height: 850 } })

	// logged-out surfaces
	await page.goto(`${BASE}/`)
	await shot(page, '01-landing-dark')
	await page.getByRole('button', { name: /show answer/i }).click()
	await shot(page, '02-landing-flashcard-flipped')
	await page.goto(`${BASE}/auth/login`)
	await shot(page, '03-login')

	// sign in as the seeded account
	await page.getByLabel('Email').fill('demo@pollyglot.dev')
	await page.getByLabel('Password', { exact: true }).fill('Password123!')
	await page.getByRole('button', { name: /sign in|login/i }).click()
	await page.waitForURL(/\/home/)
	await shot(page, '04-home-dashboard')

	await page.goto(`${BASE}/pollyglot`)
	await shot(page, '05-hub')

	// decks + deck detail + dialogs
	await page.goto(`${BASE}/pollyglot/decks`)
	await page.getByRole('link', { name: 'Japanese Basics' }).first().waitFor()
	await shot(page, '06-decks')
	await page.getByRole('link', { name: 'Japanese Basics' }).first().click()
	await page.getByRole('button', { name: /add card/i }).waitFor()
	await shot(page, '07-deck-detail')
	await page.getByRole('button', { name: /add card/i }).click()
	await page.getByLabel(/card type/i).selectOption('cloze')
	await shot(page, '08-add-card-dialog-cloze')
	await page.keyboard.press('Escape')
	await page.getByRole('button', { name: /^share$/i }).click()
	await shot(page, '09-share-dialog')
	await page.keyboard.press('Escape')

	// study session
	await page.goto(`${BASE}/pollyglot/study`)
	await shot(page, '10-study-picker')
	await page.getByRole('link', { name: /japanese basics/i }).first().click()
	const flip = page.getByRole('button', { name: /show answer/i })
	if (await flip.isVisible().catch(() => false)) {
		await shot(page, '11-study-card-front')
		await flip.click()
		await shot(page, '12-study-card-flipped')
	} else {
		await shot(page, '11-study-caught-up')
	}

	// translate with a real result
	await page.goto(`${BASE}/pollyglot/translate`)
	await page.getByLabel(/text to translate/i).fill('こんにちは')
	await page.getByRole('button', { name: 'Translate', exact: true }).click()
	await page.getByText('hello', { exact: true }).waitFor()
	await shot(page, '13-translate-result')

	// conversation with an exchange (play buttons on tutor bubbles)
	await page.goto(`${BASE}/pollyglot/conversation`)
	await shot(page, '14-conversation-list')
	await page.getByRole('button', { name: /new conversation/i }).click()
	await page.getByLabel('Language').fill('Japanese')
	await page.getByRole('button', { name: /start practicing/i }).click()
	await page.getByLabel(/tutor said/i).first().waitFor()
	await page.getByLabel('Message', { exact: true }).fill('こんにちは')
	await page.getByRole('button', { name: 'Send', exact: true }).click()
	await page.getByLabel(/you said/i).waitFor()
	await shot(page, '15-conversation-chat-with-play-buttons')

	// stats with goal progress + chart
	await page.goto(`${BASE}/pollyglot/stats`)
	await page.getByLabel(/daily goal progress/i).waitFor()
	await shot(page, '16-stats')

	// light theme spot-checks
	await page.evaluate(() => localStorage.setItem('theme', 'light'))
	await page.goto(`${BASE}/`)
	await shot(page, '17-landing-light')
	await page.goto(`${BASE}/pollyglot/decks`)
	await shot(page, '18-decks-light')
	await page.goto(`${BASE}/pollyglot/stats`)
	await page.getByLabel(/daily goal progress/i).waitFor()
	await shot(page, '19-stats-light')

	await browser.close()
	console.log('tour complete')
}

main().catch((err) => {
	console.error(err)
	process.exit(1)
})
