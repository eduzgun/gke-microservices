import type { Philosopher } from "./types"


const API_BASE = import.meta.env.GO_API_BASE || '';

const PATHS = {
  philosophers: '/philosophers',
};

export const loadPhilosophers = async (): Promise<Philosopher[]> => {
  const url = `${API_BASE}${PATHS.philosophers}`;

  try {
    const res = await fetch(url);

    // Check HTTP response
    if (!res.ok) {
      console.warn(`Failed to fetch philosophers: ${res.status} ${res.statusText}`);
      return [];
    }

    const data = await res.json();
    return data as Philosopher[];

  } catch (error) {
    console.error('Network or fetch error:', error);
    return [];
  }
};