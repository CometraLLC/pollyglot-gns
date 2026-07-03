// Speech for pronunciation and tutor messages.
// speak(): browser SpeechSynthesis, feature-detected (Pollyglot#23).
// speakWithFallback(): server audio (ElevenLabs) first, browser TTS on
// any failure — including the 503 when no provider is configured
// (Pollyglot#28).
import { speechService } from '@/src/domain/services/speech.service'

const tags: Record<string, string> = {
	japanese: 'ja-JP',
	english: 'en-US',
	spanish: 'es-ES',
	french: 'fr-FR',
	german: 'de-DE',
	italian: 'it-IT',
	portuguese: 'pt-PT',
	korean: 'ko-KR',
	chinese: 'zh-CN',
	indonesian: 'id-ID',
}

/** BCP-47 tag for a human language name; undefined lets the browser pick. */
export function languageTag(language: string): string | undefined {
	return tags[language.trim().toLowerCase()]
}

export function canSpeak(): boolean {
	return typeof window !== 'undefined' && typeof window.speechSynthesis !== 'undefined'
}

export function speak(text: string, language: string): void {
	if (!canSpeak()) return
	window.speechSynthesis.cancel()
	const utterance = new SpeechSynthesisUtterance(text)
	const tag = languageTag(language)
	if (tag) utterance.lang = tag
	window.speechSynthesis.speak(utterance)
}

export async function speakWithFallback(text: string, language: string): Promise<void> {
	try {
		const audio = await speechService.synthesize(text, language)
		const url = URL.createObjectURL(audio)
		const player = new Audio(url)
		player.onended = () => URL.revokeObjectURL(url)
		await player.play()
	} catch {
		speak(text, language)
	}
}
