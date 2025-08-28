// src/routes/philosopher/[id]/+page.server.ts
import type { PageServerLoad } from './$types';
import { apiClient } from '$lib/api/client';

export const load: PageServerLoad = async ({ fetch, params }) => {
  try {
    const philosopher = await apiClient.get(`/philosophers/${params.id}`, { fetch });
    return { philosopher };
  } catch (err) {
    return { 
      philosopher: null, 
      error: err instanceof Error ? err.message : 'Failed to load philosopher'
    };
  }
};