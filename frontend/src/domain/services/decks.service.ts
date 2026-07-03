import apiClient from './api-client';

export interface Deck {
  id: string;
  name: string;
  source_lang: string;
  target_lang: string;
  card_count: number;
  due_count: number;
  share_code?: string | null;
  created_at: string;
  updated_at: string;
}

export type CardType = 'basic' | 'cloze';

export interface Card {
  id: string;
  deck_id: string;
  front: string;
  back: string;
  card_type: CardType;
  ease_factor: number;
  interval_days: number;
  repetitions: number;
  due_at: string;
  created_at: string;
  updated_at: string;
}

export interface DeckInput {
  name: string;
  source_lang: string;
  target_lang: string;
}

export interface CardInput {
  front: string;
  back: string;
  card_type?: CardType;
  reverse?: boolean;
}

/** Study rating: 0 Forgot … 4 Got it! */
export type StudyRating = 0 | 1 | 2 | 3 | 4;

export interface SharedDeckPreview {
  name: string;
  source_lang: string;
  target_lang: string;
  card_count: number;
  sample_cards: Array<{ front: string; back: string }>;
}

export interface ImportResult {
  imported: number;
  skipped: Array<{ line: number; error: string }>;
}

export const decksService = {
  async listDecks(): Promise<Deck[]> {
    const response = await apiClient.get<Deck[]>('/v1/decks');
    return response.data;
  },

  async getDeck(id: string): Promise<Deck> {
    const response = await apiClient.get<Deck>(`/v1/decks/${id}`);
    return response.data;
  },

  async createDeck(data: DeckInput): Promise<Deck> {
    const response = await apiClient.post<Deck>('/v1/decks', data);
    return response.data;
  },

  async updateDeck(id: string, data: DeckInput): Promise<Deck> {
    const response = await apiClient.put<Deck>(`/v1/decks/${id}`, data);
    return response.data;
  },

  async deleteDeck(id: string): Promise<void> {
    await apiClient.delete(`/v1/decks/${id}`);
  },

  async listCards(deckId: string): Promise<Card[]> {
    const response = await apiClient.get<Card[]>(`/v1/decks/${deckId}/cards`);
    return response.data;
  },

  async createCard(deckId: string, data: CardInput): Promise<Card> {
    const response = await apiClient.post<Card>(`/v1/decks/${deckId}/cards`, data);
    return response.data;
  },

  async updateCard(cardId: string, data: CardInput): Promise<Card> {
    const response = await apiClient.put<Card>(`/v1/cards/${cardId}`, data);
    return response.data;
  },

  async deleteCard(cardId: string): Promise<void> {
    await apiClient.delete(`/v1/cards/${cardId}`);
  },

  async getStudyQueue(deckId: string, limit?: number): Promise<Card[]> {
    const response = await apiClient.get<Card[]>(`/v1/decks/${deckId}/queue`, {
      params: limit ? { limit } : undefined,
    });
    return response.data;
  },

  async reviewCard(cardId: string, rating: StudyRating): Promise<Card> {
    const response = await apiClient.post<Card>(`/v1/cards/${cardId}/review`, { rating });
    return response.data;
  },

  async exportDeck(deckId: string, format: 'csv' | 'tsv'): Promise<Blob> {
    const response = await apiClient.get<Blob>(`/v1/decks/${deckId}/export`, {
      params: { format },
      responseType: 'blob',
    });
    return response.data;
  },

  async shareDeck(deckId: string): Promise<{ share_code: string }> {
    const response = await apiClient.post<{ share_code: string }>(`/v1/decks/${deckId}/share`);
    return response.data;
  },

  async unshareDeck(deckId: string): Promise<void> {
    await apiClient.delete(`/v1/decks/${deckId}/share`);
  },

  async getSharedDeck(code: string): Promise<SharedDeckPreview> {
    const response = await apiClient.get<SharedDeckPreview>(`/v1/shared/${code}`);
    return response.data;
  },

  async cloneSharedDeck(code: string): Promise<Deck> {
    const response = await apiClient.post<Deck>(`/v1/shared/${code}/clone`);
    return response.data;
  },

  async importDeck(deckId: string, file: File): Promise<ImportResult> {
    const form = new FormData();
    form.append('file', file);
    const response = await apiClient.post<ImportResult>(`/v1/decks/${deckId}/import`, form, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },
};
