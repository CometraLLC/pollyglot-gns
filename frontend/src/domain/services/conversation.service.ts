import apiClient from './api-client';

export interface Conversation {
  id: string;
  title: string;
  language: string;
  created_at: string;
  updated_at: string;
}

export interface Message {
  id: string;
  role: 'user' | 'tutor';
  content: string;
  created_at: string;
}

export interface Exchange {
  user_message: Message;
  tutor_message: Message;
}

export const conversationService = {
  async listConversations(): Promise<Conversation[]> {
    const response = await apiClient.get<Conversation[]>('/v1/conversations');
    return response.data;
  },

  async createConversation(data: { title?: string; language: string }): Promise<Conversation> {
    const response = await apiClient.post<Conversation>('/v1/conversations', data);
    return response.data;
  },

  async listMessages(conversationId: string): Promise<Message[]> {
    const response = await apiClient.get<Message[]>(`/v1/conversations/${conversationId}/messages`);
    return response.data;
  },

  async sendMessage(conversationId: string, content: string): Promise<Exchange> {
    const response = await apiClient.post<Exchange>(
      `/v1/conversations/${conversationId}/messages`,
      { content }
    );
    return response.data;
  },
};
