import type { PageServerLoad } from './$types';
import { philosopherApi } from '$lib/api/philosophers';
import { interactionApi } from '$lib/api/interactions';

export const load: PageServerLoad = async ({ fetch, params }) => {
  const id = params.id;
  
  try {
    const [philosopher, interactions] = await Promise.all([
      philosopherApi.getById(id, { fetch }),
      interactionApi.get(id, { fetch })
    ]);

    return {
      philosopher,
      interactions,
    };
  } catch (err) {
    console.error('CRITICAL LOAD FAILURE:', err);
    return {
      philosopher: null,
      interactions: {
        comments: [],
        like_count: 0,
        user_liked: false
      },
      error: (err as Error).message
    };
  }
};