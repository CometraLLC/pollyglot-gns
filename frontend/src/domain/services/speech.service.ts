import apiClient from './api-client';

export const speechService = {
  /** Synthesize audio for text; throws when no provider is configured (503). */
  async synthesize(text: string, language: string): Promise<Blob> {
    const response = await apiClient.post<Blob>(
      '/v1/speech',
      { text, language },
      { responseType: 'blob' }
    );
    return response.data;
  },
};
