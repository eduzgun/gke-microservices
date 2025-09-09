import type { Philosopher } from '$lib/types';
import { apiClient } from './client';

export const philosopherApi = {
  getAll: ({ fetch }: { fetch?: typeof globalThis.fetch } = {}) =>
    apiClient.get<Philosopher[]>('/philosophers', { fetch }),

  getById: (id: number | string, { fetch }: { fetch?: typeof globalThis.fetch } = {}) =>
    apiClient.get<Philosopher>(`/philosophers/${id}`, { fetch }),

  create: (
    philosopher: Omit<Philosopher, 'id' | 'created_at'>,
    { fetch }: { fetch?: typeof globalThis.fetch } = {}
  ) =>
    apiClient.post<Philosopher>('/philosophers/add', philosopher, { fetch }),
};