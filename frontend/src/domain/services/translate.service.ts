import apiClient from './api-client';

export interface TranslateInput {
  text: string;
  from: string;
  to: string;
}

export interface Translation {
  text: string;
  from: string;
  to: string;
  translation: string;
}

export const translateService = {
  async translate(data: TranslateInput): Promise<Translation> {
    const response = await apiClient.post<Translation>('/v1/translate', data);
    return response.data;
  },
};
