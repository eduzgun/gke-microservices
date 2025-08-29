import { apiClient } from './client';
import type { Interaction } from '$lib/types';

export const interactionApi = {
    get: (
        philosopherId: string | number, 
        { fetch }: { fetch?: typeof globalThis.fetch } = {}
    ) => apiClient.get<Interaction[]>(`/philosophers/${philosopherId}/interactions`, { fetch }),


    addComment(philosopherId: number, content: string): Promise<Interaction> {
        return apiClient.post<Interaction>(
            `/philosophers/${philosopherId}/comments`,
            { content }
        );
    },

    toggleLike(philosopherId: number | string): Promise<{ liked: boolean }> {
        return apiClient.post<{ liked: boolean }>(
            `/philosophers/${philosopherId}/like`,
            {}
        );
    }
};