import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  decksService,
  DeckInput,
  CardInput,
  StudyRating,
} from '@/src/domain/services/decks.service';

// Query keys
export const deckKeys = {
  all: ['decks'] as const,
  lists: () => [...deckKeys.all, 'list'] as const,
  detail: (id: string) => [...deckKeys.all, 'detail', id] as const,
  cards: (deckId: string) => [...deckKeys.all, 'cards', deckId] as const,
  queue: (deckId: string) => [...deckKeys.all, 'queue', deckId] as const,
};

export function useDecks() {
  return useQuery({
    queryKey: deckKeys.lists(),
    queryFn: () => decksService.listDecks(),
  });
}

export function useDeck(id: string) {
  return useQuery({
    queryKey: deckKeys.detail(id),
    queryFn: () => decksService.getDeck(id),
    enabled: !!id,
  });
}

export function useCreateDeck() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: DeckInput) => decksService.createDeck(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deckKeys.lists() });
    },
  });
}

export function useUpdateDeck() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: DeckInput }) =>
      decksService.updateDeck(id, data),
    onSuccess: (_deck, { id }) => {
      queryClient.invalidateQueries({ queryKey: deckKeys.lists() });
      queryClient.invalidateQueries({ queryKey: deckKeys.detail(id) });
    },
  });
}

export function useDeleteDeck() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => decksService.deleteDeck(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deckKeys.lists() });
    },
  });
}

export function useCards(deckId: string) {
  return useQuery({
    queryKey: deckKeys.cards(deckId),
    queryFn: () => decksService.listCards(deckId),
    enabled: !!deckId,
  });
}

export function useCreateCard(deckId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CardInput) => decksService.createCard(deckId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deckKeys.cards(deckId) });
      queryClient.invalidateQueries({ queryKey: deckKeys.detail(deckId) });
      queryClient.invalidateQueries({ queryKey: deckKeys.lists() });
    },
  });
}

export function useUpdateCard(deckId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ cardId, data }: { cardId: string; data: CardInput }) =>
      decksService.updateCard(cardId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deckKeys.cards(deckId) });
    },
  });
}

export function useDeleteCard(deckId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (cardId: string) => decksService.deleteCard(cardId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deckKeys.cards(deckId) });
      queryClient.invalidateQueries({ queryKey: deckKeys.detail(deckId) });
      queryClient.invalidateQueries({ queryKey: deckKeys.lists() });
    },
  });
}

export function useStudyQueue(deckId: string, limit?: number) {
  return useQuery({
    queryKey: deckKeys.queue(deckId),
    queryFn: () => decksService.getStudyQueue(deckId, limit),
    enabled: !!deckId,
    // The queue is consumed as the user studies; never serve a stale one.
    staleTime: 0,
  });
}

export function useReviewCard(deckId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ cardId, rating }: { cardId: string; rating: StudyRating }) =>
      decksService.reviewCard(cardId, rating),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deckKeys.cards(deckId) });
    },
  });
}
