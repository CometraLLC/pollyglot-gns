import { afterEach, describe, expect, it, vi } from 'vitest'
import { canSpeak, languageTag, speak } from './speech'

afterEach(() => {
	vi.unstubAllGlobals()
})

describe('languageTag', () => {
	it('maps language names to BCP-47 tags', () => {
		expect(languageTag('Japanese')).toBe('ja-JP')
		expect(languageTag('japanese')).toBe('ja-JP')
		expect(languageTag('Spanish')).toBe('es-ES')
		expect(languageTag('English')).toBe('en-US')
		expect(languageTag('German')).toBe('de-DE')
		expect(languageTag('French')).toBe('fr-FR')
	})

	it('falls back to undefined for unknown languages (browser default voice)', () => {
		expect(languageTag('Klingon')).toBeUndefined()
	})
})

describe('speak', () => {
	it('is unsupported when the browser lacks speechSynthesis', () => {
		vi.stubGlobal('speechSynthesis', undefined)

		expect(canSpeak()).toBe(false)
	})

	it('speaks with the mapped language', () => {
		const speakMock = vi.fn()
		vi.stubGlobal('speechSynthesis', { speak: speakMock, cancel: vi.fn() })
		vi.stubGlobal(
			'SpeechSynthesisUtterance',
			class {
				text: string
				lang = ''
				constructor(text: string) {
					this.text = text
				}
			}
		)

		expect(canSpeak()).toBe(true)
		speak('こんにちは', 'Japanese')

		expect(speakMock).toHaveBeenCalledTimes(1)
		const utterance = speakMock.mock.calls[0][0]
		expect(utterance.text).toBe('こんにちは')
		expect(utterance.lang).toBe('ja-JP')
	})
})
