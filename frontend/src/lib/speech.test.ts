import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('@/src/domain/services/speech.service', () => ({
	speechService: {
		synthesize: vi.fn(),
	},
}))

import { speechService } from '@/src/domain/services/speech.service'
import { canSpeak, languageTag, speak, speakWithFallback } from './speech'

const mockedSpeech = vi.mocked(speechService)

beforeEach(() => {
	vi.clearAllMocks()
})

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

describe('speakWithFallback', () => {
	it('plays server audio when the provider is configured', async () => {
		const play = vi.fn().mockResolvedValue(undefined)
		vi.stubGlobal('Audio', class { onended: (() => void) | null = null; play = play })
		vi.stubGlobal('URL', {
			createObjectURL: vi.fn(() => 'blob:audio'),
			revokeObjectURL: vi.fn(),
		})
		mockedSpeech.synthesize.mockResolvedValue(new Blob(['mp3']))

		await speakWithFallback('こんにちは', 'Japanese')

		expect(mockedSpeech.synthesize).toHaveBeenCalledWith('こんにちは', 'Japanese')
		expect(play).toHaveBeenCalledTimes(1)
	})

	it('falls back to browser TTS when the server declines', async () => {
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
		mockedSpeech.synthesize.mockRejectedValue(new Error('503 not configured'))

		await speakWithFallback('こんにちは', 'Japanese')

		expect(speakMock).toHaveBeenCalledTimes(1)
		expect(speakMock.mock.calls[0][0].text).toBe('こんにちは')
	})
})
