import type { PageServerLoad } from './$types';
import { authApi } from '$lib/api/auth';

export const load: PageServerLoad = async ({ fetch }) => {
  try {
    // This will send cookies automatically in server-side fetch
    const response = await fetch('/profile');
    if (response.ok) {
      const user = await response.json();
      return { user };
    }
    return { user: null };
  } catch {
    return { user: null };
  }
};