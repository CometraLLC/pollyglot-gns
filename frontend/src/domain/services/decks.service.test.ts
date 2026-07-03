import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('./api-client', () => ({
	default: {
		get: vi.fn(),
		post: vi.fn(),
		put: vi.fn(),
		delete: vi.fn(),
	},
}))

import apiClient from './api-client'
import { decksService } from './decks.service'

const mocked = vi.mocked(apiClient)

beforeEach(() => {
	vi.clearAllMocks()
	mocked.get.mockResolvedValue({ data: [] })
	mocked.post.mockResolvedValue({ data: {} })
	mocked.put.mockResolvedValue({ data: {} })
	mocked.delete.mockResolvedValue({ data: {} })
})

// The API base URL is versionless; every path must carry /v1 exactly once.
// This suite pins the full request contract (paths, verbs, payloads) so a
// regression like the /v1/v1 double-prefix bug (Pollyglot#19) fails fast.
describe('decksService request contract', () => {
	it('lists decks via GET /v1/decks', async () => {
		await decksService.listDecks()
		expect(mocked.get).toHaveBeenCalledWith('/v1/decks')
	})

	it('gets a deck via GET /v1/decks/{id}', async () => {
		await decksService.getDeck('d1')
		expect(mocked.get).toHaveBeenCalledWith('/v1/decks/d1')
	})

	it('creates a deck via POST /v1/decks with the exact payload', async () => {
		const input = { name: 'N', source_lang: 'Japanese', target_lang: 'English' }
		await decksService.createDeck(input)
		expect(mocked.post).toHaveBeenCalledWith('/v1/decks', input)
	})

	it('updates a deck via PUT /v1/decks/{id}', async () => {
		const input = { name: 'N2', source_lang: 'a', target_lang: 'b' }
		await decksService.updateDeck('d1', input)
		expect(mocked.put).toHaveBeenCalledWith('/v1/decks/d1', input)
	})

	it('deletes a deck via DELETE /v1/decks/{id}', async () => {
		await decksService.deleteDeck('d1')
		expect(mocked.delete).toHaveBeenCalledWith('/v1/decks/d1')
	})

	it('lists cards via GET /v1/decks/{id}/cards', async () => {
		await decksService.listCards('d1')
		expect(mocked.get).toHaveBeenCalledWith('/v1/decks/d1/cards')
	})

	it('creates a card via POST /v1/decks/{id}/cards', async () => {
		await decksService.createCard('d1', { front: 'f', back: 'b' })
		expect(mocked.post).toHaveBeenCalledWith('/v1/decks/d1/cards', { front: 'f', back: 'b' })
	})

	it('passes card_type and reverse through on create', async () => {
		await decksService.createCard('d1', {
			front: '{{c1::猫}}が好き',
			back: 'b',
			card_type: 'cloze',
		})
		expect(mocked.post).toHaveBeenCalledWith('/v1/decks/d1/cards', {
			front: '{{c1::猫}}が好き',
			back: 'b',
			card_type: 'cloze',
		})

		await decksService.createCard('d1', { front: 'f', back: 'b', reverse: true })
		expect(mocked.post).toHaveBeenCalledWith('/v1/decks/d1/cards', {
			front: 'f',
			back: 'b',
			reverse: true,
		})
	})

	it('updates a card via PUT /v1/cards/{id}', async () => {
		await decksService.updateCard('c1', { front: 'f', back: 'b' })
		expect(mocked.put).toHaveBeenCalledWith('/v1/cards/c1', { front: 'f', back: 'b' })
	})

	it('deletes a card via DELETE /v1/cards/{id}', async () => {
		await decksService.deleteCard('c1')
		expect(mocked.delete).toHaveBeenCalledWith('/v1/cards/c1')
	})

	it('fetches the study queue via GET /v1/decks/{id}/queue with optional limit', async () => {
		await decksService.getStudyQueue('d1')
		expect(mocked.get).toHaveBeenCalledWith('/v1/decks/d1/queue', { params: undefined })

		await decksService.getStudyQueue('d1', 10)
		expect(mocked.get).toHaveBeenCalledWith('/v1/decks/d1/queue', { params: { limit: 10 } })
	})

	it('reviews a card via POST /v1/cards/{id}/review with the rating', async () => {
		await decksService.reviewCard('c1', 0)
		expect(mocked.post).toHaveBeenCalledWith('/v1/cards/c1/review', { rating: 0 })

		await decksService.reviewCard('c1', 4)
		expect(mocked.post).toHaveBeenCalledWith('/v1/cards/c1/review', { rating: 4 })
	})

	it('exports a deck via GET /v1/decks/{id}/export expecting a blob', async () => {
		mocked.get.mockResolvedValue({ data: new Blob(['csv']) })

		await decksService.exportDeck('d1', 'tsv')

		expect(mocked.get).toHaveBeenCalledWith('/v1/decks/d1/export', {
			params: { format: 'tsv' },
			responseType: 'blob',
		})
	})

	it('imports a deck via multipart POST /v1/decks/{id}/import', async () => {
		mocked.post.mockResolvedValue({ data: { imported: 2, skipped: [] } })
		const file = new File(['a,b'], 'cards.csv', { type: 'text/csv' })

		const result = await decksService.importDeck('d1', file)

		expect(mocked.post).toHaveBeenCalledWith(
			'/v1/decks/d1/import',
			expect.any(FormData),
			{ headers: { 'Content-Type': 'multipart/form-data' } }
		)
		const sentForm = mocked.post.mock.calls[0][1] as FormData
		expect(sentForm.get('file')).toBe(file)
		expect(result.imported).toBe(2)
	})

	it('shares, unshares, previews, and clones via the shared endpoints', async () => {
		mocked.post.mockResolvedValue({ data: { share_code: 'ABCDEF2345' } })
		await decksService.shareDeck('d1')
		expect(mocked.post).toHaveBeenCalledWith('/v1/decks/d1/share')

		await decksService.unshareDeck('d1')
		expect(mocked.delete).toHaveBeenCalledWith('/v1/decks/d1/share')

		mocked.get.mockResolvedValue({ data: { name: 'X' } })
		await decksService.getSharedDeck('ABCDEF2345')
		expect(mocked.get).toHaveBeenCalledWith('/v1/shared/ABCDEF2345')

		await decksService.cloneSharedDeck('ABCDEF2345')
		expect(mocked.post).toHaveBeenCalledWith('/v1/shared/ABCDEF2345/clone')
	})

	it('unwraps response data', async () => {
		const decks = [{ id: 'd1' }]
		mocked.get.mockResolvedValue({ data: decks })
		await expect(decksService.listDecks()).resolves.toEqual(decks)
	})
})
