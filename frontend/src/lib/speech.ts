// Browser speech synthesis for pronunciation (Pollyglot#23).
// Feature-detected: callers hide their UI when canSpeak() is false.

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
