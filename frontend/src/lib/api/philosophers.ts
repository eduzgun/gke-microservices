import type { Philosopher } from '$lib/types';
import { apiClient } from './client';

export const philosopherApi = {
  getAll: () => apiClient.get<Philosopher[]>('/philosophers'),
  getById: (id: number | string) => apiClient.get<Philosopher>(`/philosophers/${id}`),
  create: (philosopher: Omit<Philosopher, 'id' | 'created_at'>) =>
    apiClient.post<Philosopher>('/philosophers/add', philosopher),
};