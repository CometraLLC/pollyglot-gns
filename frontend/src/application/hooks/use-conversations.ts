import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  conversationService,
  Exchange,
  Message,
} from '@/src/domain/services/conversation.service';

export const conversationKeys = {
  all: ['conversations'] as const,
  lists: () => [...conversationKeys.all, 'list'] as const,
  messages: (conversationId: string) =>
    [...conversationKeys.all, 'messages', conversationId] as const,
};

export function useConversations() {
  return useQuery({
    queryKey: conversationKeys.lists(),
    queryFn: () => conversationService.listConversations(),
  });
}

export function useCreateConversation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { title?: string; language: string }) =>
      conversationService.createConversation(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: conversationKeys.lists() });
    },
  });
}

export function useMessages(conversationId: string) {
  return useQuery({
    queryKey: conversationKeys.messages(conversationId),
    queryFn: () => conversationService.listMessages(conversationId),
    enabled: !!conversationId,
  });
}

export function useSendMessage(conversationId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (content: string) => conversationService.sendMessage(conversationId, content),
    onSuccess: (exchange: Exchange) => {
      // Append the exchange locally — no refetch, so the chat never flickers.
      queryClient.setQueryData<Message[]>(
        conversationKeys.messages(conversationId),
        (existing = []) => [...existing, exchange.user_message, exchange.tutor_message]
      );
    },
  });
}
