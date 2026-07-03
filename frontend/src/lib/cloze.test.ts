import { describe, expect, it } from 'vitest'
import { blank, deletions, reveal } from './cloze'

// Mirrors backend/pkg/cloze — keep the two in sync.
describe('cloze', () => {
	it('extracts deletions in order', () => {
		expect(deletions('水を{{c1::飲みます}}')).toEqual(['飲みます'])
		expect(deletions('{{c1::猫}}が{{c2::魚}}を食べる')).toEqual(['猫', '魚'])
	})

	it('treats malformed markers as plain text', () => {
		expect(deletions('plain text')).toEqual([])
		expect(deletions('水を{{c1::飲みます')).toEqual([])
		expect(deletions('{{飲みます}}')).toEqual([])
		expect(deletions('水を{{c1::}}')).toEqual([])
	})

	it('blanks deletions for the card front', () => {
		expect(blank('水を{{c1::飲みます}}')).toBe('水を[…]')
		expect(blank('{{c1::猫}}が{{c2::魚}}を食べる')).toBe('[…]が[…]を食べる')
		expect(blank('plain text')).toBe('plain text')
	})

	it('reveals the full text for the flip', () => {
		expect(reveal('水を{{c1::飲みます}}')).toBe('水を飲みます')
		expect(reveal('{{c1::猫}}が{{c2::魚}}を食べる')).toBe('猫が魚を食べる')
	})
})
