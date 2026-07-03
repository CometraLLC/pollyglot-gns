import apiClient from './api-client';

export interface DayCount {
  date: string;
  count: number;
}

export interface Stats {
  reviews_today: number;
  total_reviews: number;
  unique_cards: number;
  streak_days: number;
  reviews_per_day: DayCount[];
}

export const statsService = {
  async getStats(): Promise<Stats> {
    const response = await apiClient.get<Stats>('/v1/stats');
    return response.data;
  },
};
