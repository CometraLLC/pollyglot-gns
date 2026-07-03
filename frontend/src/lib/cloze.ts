// Anki-style cloze deletions: {{c1::hidden text}}.
// Mirrors backend/pkg/cloze — keep the two in sync.
const marker = /\{\{c\d+::(.+?)\}\}/g

export function deletions(front: string): string[] {
	return Array.from(front.matchAll(marker), (m) => m[1])
}

export function blank(front: string): string {
	return front.replace(marker, '[…]')
}

export function reveal(front: string): string {
	return front.replace(marker, '$1')
}
