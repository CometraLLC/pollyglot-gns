import { describe, expect, it } from 'vitest'
import en from '@/locales/en.json'
import id from '@/locales/id.json'

// Every key must exist in both languages — a missing translation fails
// here instead of rendering a raw key in production.
function keyPaths(obj: Record<string, unknown>, prefix = ''): string[] {
	return Object.entries(obj).flatMap(([key, value]) => {
		const path = prefix ? `${prefix}.${key}` : key
		if (value && typeof value === 'object' && !Array.isArray(value)) {
			return keyPaths(value as Record<string, unknown>, path)
		}
		return [path]
	})
}

describe('locale catalogs', () => {
	it('en and id define exactly the same keys', () => {
		const enKeys = keyPaths(en).sort()
		const idKeys = keyPaths(id).sort()

		expect(idKeys).toEqual(enKeys)
	})

	it('covers the Pollyglot navigation', () => {
		const keys = keyPaths(en)

		for (const key of [
			'nav.pollyglot',
			'nav.decks',
			'nav.study',
			'nav.translate',
			'nav.conversation',
			'nav.progress',
		]) {
			expect(keys).toContain(key)
		}
	})
})
